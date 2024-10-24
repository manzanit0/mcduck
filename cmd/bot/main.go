package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/otelconnect"
	"github.com/gin-gonic/gin"
	"github.com/manzanit0/mcduck/cmd/bot/internal/bot"
	"github.com/manzanit0/mcduck/gen/api/receipts.v1/receiptsv1connect"
	"github.com/manzanit0/mcduck/gen/api/users.v1/usersv1connect"
	"github.com/manzanit0/mcduck/internal/mcduck"
	"github.com/manzanit0/mcduck/pkg/micro"
	"github.com/manzanit0/mcduck/pkg/tgram"
	"github.com/manzanit0/mcduck/pkg/xhttp"
	"github.com/manzanit0/mcduck/pkg/xtrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	serviceName = "tgram-bot"
)

func main() {
	svc, err := micro.NewGinService(serviceName)
	if err != nil {
		panic(err)
	}

	interceptor, _ := otelconnect.NewInterceptor()
	receiptsClient := receiptsv1connect.NewReceiptsServiceClient(xhttp.NewClient(), micro.MustGetEnv("PRIVATE_DOTS_HOST"), connect.WithInterceptors(interceptor))
	usersClient := usersv1connect.NewUsersServiceClient(xhttp.NewClient(), micro.MustGetEnv("PRIVATE_DOTS_HOST"), connect.WithInterceptors(interceptor))
	uploader := mcduck.NewReceiptUploader(usersClient, receiptsClient)

	tgramClient := tgram.NewClient(xhttp.NewClient(), micro.MustGetEnv("TELEGRAM_BOT_TOKEN"))

	svc.Engine.POST("/telegram/webhook", telegramWebhookController(tgramClient, uploader))

	if err := svc.Run(); err != nil {
		slog.Error("run ended with error", "error", err.Error())
		os.Exit(1)
	}
}

func telegramWebhookController(tgramClient tgram.Client, uploader mcduck.ReceiptUploader) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, span := xtrace.GetSpan(c.Request.Context())

		var r tgram.WebhookRequest
		if err := c.ShouldBindJSON(&r); err != nil {
			span.SetStatus(codes.Error, err.Error())
			slog.ErrorContext(ctx, "unable to bind json", "error", err.Error())

			// FIXME: this actually isn't what's happening. It's not a json.Unmarshall as I expected.
			res := gin.H{"error": fmt.Sprintf("payload does not conform with telegram contract: %s", err.Error())}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		span.SetAttributes(
			attribute.Int64("mduck.telegram.chat_id", r.GetChatID()),
			attribute.String("mduck.telegram.language_code", r.GetFromLanguageCode()),
		)

		switch {
		case r.Message != nil && r.Message.Text != nil && strings.HasPrefix(*r.Message.Text, "/login"):
			span.SetAttributes(attribute.String("mduck.telegram.command", "login"))

			res := bot.LoginLink(ctx, &r)
			c.JSON(http.StatusOK, res)

			// The message has either photos or a doc.
		case r.Message != nil && (len(r.Message.Photos) > 0 || r.Message.Document != nil):
			span.SetAttributes(attribute.String("mduck.telegram.command", "upload"))

			res := bot.UploadReceipt(ctx, tgramClient, uploader, &r)
			c.JSON(http.StatusOK, res)

		default:
			span.SetAttributes(attribute.String("mduck.telegram.command", "unknown"))

			// NOTE: If it's the message is sent in a group, we don't want to spam
			// the group with "Hey!" messages.
			if r.GetChatID() != r.GetFromID() {
				c.JSON(http.StatusOK, "")
				return
			}

			res := tgram.NewMarkdownResponse("Hey\\! Just send me a picture with a receipt and I will do the rest\\!", r.GetChatID())
			c.JSON(http.StatusOK, res)
		}
	}
}
