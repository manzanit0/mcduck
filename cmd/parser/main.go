package main

import (
	"context"
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
	"github.com/manzanit0/mcduck/internal/mcduck"
	"github.com/manzanit0/mcduck/internal/parser"
	"github.com/manzanit0/mcduck/pkg/micro"
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

	// NOTE: AWS SDK expects the environment variables, so assert their existence.
	_ = micro.MustGetEnv("AWS_ACCESS_KEY")
	_ = micro.MustGetEnv("AWS_SECRET_ACCESS_KEY")

	apiKey := micro.MustGetEnv("OPENAI_API_KEY")
	config, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		return err
	}

	pdfParser := parser.NewTextractParser(config, apiKey)
	imageParser := parser.NewAIVisionParser(apiKey)
	smartParser := parser.NewSmartParser(pdfParser, imageParser)

	// NOTE: this is curently not being used. That's ok though. It's handy to
	// test behaviour manually since the pubsub mechanism relies on proto
	// encoding.
	r.POST("/receipt", endpoint(smartParser))

	natsURL := micro.MustGetEnv("NATS_URL")

	slog.InfoContext(context.Background(), "ensuring stream exists")
	_, _, err = pubsub.NewStream(context.Background(), natsURL, pubsub.DefaultStreamName, "events.receipts.v1.ReceiptCreated")
	if err != nil {
		return err
	}

	consumer := func(ctx context.Context) error {
		slog.InfoContext(ctx, "starting consumer in the background")
		err = ConsumeMessages(ctx, natsURL, dbx, smartParser)
		if err != nil && err != context.Canceled {
			return err
		}

		return nil
	}

	return micro.RunGracefully(r, consumer)
}

func endpoint(parser parser.ReceiptParser) func(c *gin.Context) {
	return func(c *gin.Context) {
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
			slog.ErrorContext(ctx, "unable to save file to disk", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to save file to disk: %s", err.Error())})
			return
		}

		span.SetAttributes(attribute.String("file.path", filePath))

		data, err := os.ReadFile(filePath)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			slog.ErrorContext(ctx, "unable to read file from disk", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable to read file from disk: %s", err.Error())})
			return
		}

		receipt, _, err := parser.ExtractReceipt(ctx, data)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			slog.ErrorContext(ctx, "failed to extract receipt", "error", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unable extract data from receipt: %s", err.Error())})
			return
		}

		c.JSON(http.StatusOK, receipt)
	}
}

func ConsumeMessages(ctx context.Context, natsURL string, dbx *sqlx.DB, parser parser.ReceiptParser) error {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("connect nats: %w", err)
	}

	jss, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("open jetstream: %w", err)
	}

	const maxDeliveries = 5
	const randomHighNumber = 8192
	msgCh := make(chan *nats.Msg, randomHighNumber)
	_, err = jss.ChanQueueSubscribe("events.receipts.v1.ReceiptCreated", serviceName, msgCh, nats.ManualAck(), nats.MaxDeliver(maxDeliveries))
	if err != nil {
		return fmt.Errorf("subscribe queue: %w", err)
	}

	const randomWorkerCount = 100
	const randomPoolCapacity = 1000
	pool := pond.New(randomWorkerCount, randomPoolCapacity, pond.Context(ctx))
	defer pool.StopAndWait()

	augmentor := mcduck.AIaugmentor{DB: dbx, Parser: parser}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg := <-msgCh:
			pool.Submit(func() {
				slog.InfoContext(ctx, "Received events.receipts.v1.ReceiptCreated")

				meta, err := msg.Metadata()
				if err != nil {
					slog.ErrorContext(ctx, "failed to get message metadata", "error", err.Error())
					return
				}

				switch {
				case meta.NumDelivered == maxDeliveries:
					if err = processFailedMessage(msg, dbx); err != nil {
						slog.InfoContext(ctx, "failed to mark as unprocessable")
						return
					}

					err = msg.Term()
					if err != nil {
						slog.ErrorContext(ctx, "failed to TERM message", "error", err.Error())
					}

				default:
					err = processMessage(msg, &augmentor)
					if err != nil {
						slog.ErrorContext(ctx, "failed to process message", "error", err.Error())
						err = msg.Nak()
						if err != nil {
							slog.ErrorContext(ctx, "failed to NACK message", "error", err.Error())
						}

						return
					}

					err = msg.Ack()
					if err != nil {
						slog.ErrorContext(ctx, "failed to ACK message", "error", err.Error())
					}

					slog.InfoContext(ctx, "correctly processed message")
				}
			})

		case <-time.After(1 * time.Second):
		}
	}
}

func processFailedMessage(msg *nats.Msg, dbx *sqlx.DB) error {
	ctx, span := xtrace.StartSpan(context.Background(), "Consume ReceiptCreated event")
	defer span.End()

	ev, err := pubsub.UnmarshalProto(msg.Data, &receiptsevv1.ReceiptCreated{})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	receiptRepo := mcduck.NewReceiptRepository(dbx)
	err = receiptRepo.MarkFailedToProcess(ctx, ev.Receipt.Id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

func processMessage(msg *nats.Msg, augmentor mcduck.ReceiptAugmentor) error {
	// FIXME: this trace_id/span_id context propagation isn't working.
	ctx := xtrace.HydrateContext(context.Background(), msg.Header.Get("trace_id"), msg.Header.Get("span_id"))
	ctx, span := xtrace.StartSpan(ctx, "Consume ReceiptCreated event")
	defer span.End()

	ev, err := pubsub.UnmarshalProto(msg.Data, &receiptsevv1.ReceiptCreated{})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	err = augmentor.AugmentCreatedReceipt(ctx, ev)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}
