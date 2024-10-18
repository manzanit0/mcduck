package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alitto/pond"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	receiptsevv1 "github.com/manzanit0/mcduck/gen/events/receipts.v1"
	"github.com/manzanit0/mcduck/internal/expense"
	"github.com/manzanit0/mcduck/internal/parser"
	"github.com/manzanit0/mcduck/internal/receipt"
	"github.com/manzanit0/mcduck/pkg/micro"
	"github.com/manzanit0/mcduck/pkg/openai"
	"github.com/manzanit0/mcduck/pkg/pubsub"
	"github.com/manzanit0/mcduck/pkg/xlog"
	"github.com/manzanit0/mcduck/pkg/xsql"
	"github.com/manzanit0/mcduck/pkg/xtrace"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	serviceName = "parser"
	awsRegion   = "eu-west-1"
)

func main() {
	if err := run(); err != nil {
		slog.Error("exiting server", "error", err.Error())
		os.Exit(1)
	}

	slog.Info("exiting server")
}

func run() error {
	xlog.InitSlog()

	tp, err := xtrace.TracerFromEnv(context.Background(), serviceName)
	if err != nil {
		return err
	}
	defer tp.Shutdown(context.Background())

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(xlog.EnhanceContext)
	r.Use(tp.TraceRequests())
	r.Use(tp.EnhanceTraceMetadata())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	dbx, err := xsql.OpenFromEnv()
	if err != nil {
		return err
	}
	defer xsql.Close(dbx)

	apiKey := micro.MustGetEnv("OPENAI_API_KEY")

	// Just crash the service if these aren't available.
	_ = micro.MustGetEnv("AWS_ACCESS_KEY")
	_ = micro.MustGetEnv("AWS_SECRET_ACCESS_KEY")

	config, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		return err
	}

	textractParser := parser.NewTextractParser(config, apiKey)
	aivisionParser := parser.NewAIVisionParser(apiKey)

	// NOTE: this is curently not being used. That's ok though. It's handy to
	// test behaviour manually since the pubsub mechanism relies on proto
	// encoding.
	r.POST("/receipt", func(c *gin.Context) {
		ctx, span := xtrace.GetSpan(c.Request.Context())

		file, err := c.FormFile("receipt")
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to read file: %s", err.Error())})
			return
		}

		dir := os.TempDir()
		filePath := filepath.Join(dir, file.Filename)
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to save file to disk: %s", err.Error())})
			return
		}

		span.SetAttributes(attribute.String("file.path", filePath))

		data, err := os.ReadFile(filePath)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("unable to read file from disk: %s", err.Error())})
			return
		}

		receipt, err := parseReceipt(ctx, data, textractParser, aivisionParser)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			slog.ErrorContext(c.Request.Context(), "failed to extract receipt", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable extract data from receipt: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, receipt)
	})

	natsURL := micro.MustGetEnv("NATS_URL")

	slog.InfoContext(context.Background(), "ensuring stream exists")
	_, _, err = pubsub.NewStream(context.Background(), natsURL, pubsub.DefaultStreamName, "events.receipts.v1.ReceiptCreated")
	if err != nil {
		return err
	}

	consumer := func(ctx context.Context) error {
		slog.InfoContext(ctx, "starting consumer in the background")
		err = ConsumeMessages(ctx, natsURL, dbx, textractParser, aivisionParser)
		if err != nil && err != context.Canceled {
			return err
		}

		return nil
	}

	return micro.RunGracefully(r, consumer)
}

func ConsumeMessages(ctx context.Context, natsURL string, dbx *sqlx.DB, pdfParser parser.ReceiptParser, imageParser parser.ReceiptParser) error {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("connect nats: %w", err)
	}

	jss, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("open jetstream: %w", err)
	}

	msgCh := make(chan *nats.Msg, 8192)

	// FIXME: we should actually AckTerm() the messages manually and update the
	// receipt status to "failed to process" or something similar.
	_, err = jss.ChanQueueSubscribe("events.receipts.v1.ReceiptCreated", serviceName, msgCh, nats.ManualAck(), nats.MaxDeliver(5))
	if err != nil {
		return fmt.Errorf("subscribe queue: %w", err)
	}

	// TODO: might be better to create more queue subscribers instead.
	pool := pond.New(100, 1000, pond.Context(ctx))
	defer pool.StopAndWait()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg := <-msgCh:
			pool.Submit(func() {
				slog.InfoContext(ctx, "Received events.receipts.v1.ReceiptCreated")

				err := processMessage(msg, dbx, pdfParser, imageParser)
				if err != nil {
					slog.ErrorContext(ctx, "failed to process message", "error", err.Error())
					err = msg.Nak()
					if err != nil {
						slog.ErrorContext(ctx, "failed to NACK message", "error", err.Error())
					}
				} else {
					err = msg.Ack()
					if err != nil {
						slog.ErrorContext(ctx, "failed to ACK message", "error", err.Error())
					}
				}
			})

		case <-time.After(1 * time.Second):
		}
	}
}

func processMessage(msg *nats.Msg, dbx *sqlx.DB, pdfParser parser.ReceiptParser, imageParser parser.ReceiptParser) error {
	// FIXME: this trace_id/span_id context propagation isn't working.
	ctx := xtrace.HydrateContext(context.Background(), msg.Header.Get("trace_id"), msg.Header.Get("span_id"))
	ctx, span := xtrace.StartSpan(ctx, "Consume ReceiptCreated event")
	defer span.End()

	ev, err := pubsub.UnmarshalProto(msg.Data, &receiptsevv1.ReceiptCreated{})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	data := ev.Receipt.File
	parsedReceipt, err := parseReceipt(ctx, data, pdfParser, imageParser)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	parsedTime, err := time.Parse("02/01/2006", parsedReceipt.PurchaseDate)
	if err != nil {
		slog.Info("failed to parse receipt date. Defaulting to 'now' ", "error", err.Error())
		parsedTime = time.Now()
	}

	txn, err := dbx.BeginTxx(ctx, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	defer xsql.TxClose(txn)

	receiptRepo := receipt.NewRepository(dbx)
	err = receiptRepo.UpdateReceiptWithTxn(ctx, txn, receipt.UpdateReceiptRequest{
		ID:     int64(ev.Receipt.Id),
		Vendor: &parsedReceipt.Vendor,
		Date:   &parsedTime,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = expense.CreateExpenses(ctx, txn, expense.ExpensesBatch{
		Records: []expense.Expense{
			{
				Date:        parsedTime,
				Amount:      float32(parsedReceipt.Amount),
				Category:    "Receipt Upload",
				Subcategory: "",
				UserEmail:   ev.UserEmail,
				ReceiptID:   ev.Receipt.Id,
				Description: "This expense has be autogenerated by the system",
			},
		},
		UserEmail: ev.UserEmail,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err = txn.Commit(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func parseReceipt(ctx context.Context, data []byte, pdfParser parser.ReceiptParser, imageParser parser.ReceiptParser) (receipt *parser.Receipt, err error) {
	_, span := xtrace.GetSpan(ctx)
	span.SetAttributes(attribute.Int("file.size", len(data)))

	contentType := http.DetectContentType(data)
	span.SetAttributes(attribute.String("file.content_type", contentType))

	var openAIRes *openai.Response

	switch contentType {
	case "application/pdf":
		receipt, openAIRes, err = pdfParser.ExtractReceipt(ctx, data)
		if err != nil {
			return
		}

	// Default to images
	default:
		receipt, openAIRes, err = imageParser.ExtractReceipt(ctx, data)
		if err != nil {
			return
		}
	}

	// NOTE: currently this is where I'm putting my money for checking responses
	// and evaluating LLM perf. Not ideal, but good enough for now.
	marshalled, _ := json.Marshal(receipt)
	marshalledRes, _ := json.Marshal(openAIRes)
	span.SetAttributes(attribute.String("openai.response", string(marshalledRes)))
	slog.InfoContext(ctx, "chatGPT response", "processed_receipt", marshalled, "open_ai_response", marshalledRes)

	return
}
