package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/manzanit0/mcduck/internal/mcduck"
	"github.com/manzanit0/mcduck/pkg/tgram"
	"github.com/manzanit0/mcduck/pkg/xtrace"
	"go.opentelemetry.io/otel/codes"
)

func GetDocument(ctx context.Context, tgramClient tgram.Client, fileID string) ([]byte, error) {
	_, span := xtrace.StartSpan(ctx, "telegram.GetFile")
	file, err := tgramClient.GetFile(tgram.GetFileRequest{FileID: fileID})
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("get file: %w", err)
	}
	span.End()

	_, span = xtrace.StartSpan(ctx, "telegram.DownloadFile")
	fileData, err := tgramClient.DownloadFile(file)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("download file: %w", err)
	}
	span.End()

	return fileData, nil
}

func UploadReceipt(ctx context.Context, tgramClient tgram.Client, uploader mcduck.ReceiptUploader, r *tgram.WebhookRequest) *tgram.WebhookResponse {
	ctx, span := xtrace.StartSpan(ctx, "Upload Receipt via Telegram")
	defer span.End()

	var fileID string
	var fileSize int64

	if r.Message.Document != nil {
		fileID = r.Message.Document.FileID
	} else if len(r.Message.Photos) > 0 {
		// Get the biggest photo: this will ensure better parsing by parser service.
		for _, p := range r.Message.Photos {
			if p.FileSize != nil && *p.FileSize > fileSize {
				fileID = p.FileID
				fileSize = *p.FileSize
			}
		}
	}

	fileData, err := GetDocument(ctx, tgramClient, fileID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "tgram.DownloadFile:", "error", err.Error())
		return tgram.NewHTMLResponse(fmt.Sprintf("unable to download file from Telegram servers: %s", err.Error()), r.GetChatID())
	}

	if len(fileData) == 0 {
		return tgram.NewHTMLResponse("empty file", r.GetChatID())
	}

	id, err := uploader.UploadFromChat(ctx, fileData, r.GetChatID())
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "Upload receipt", "error", err.Error())
		return tgram.NewHTMLResponse(fmt.Sprintf("unable to upload receipt: %s", err.Error()), r.GetChatID())
	}

	return tgram.NewHTMLResponse(fmt.Sprintf("Receipt submitted for processing. It's ID is %d.", id), r.GetChatID())
}
