package mcduck

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	receiptsv1 "github.com/manzanit0/mcduck/gen/api/receipts.v1"
	"github.com/manzanit0/mcduck/gen/api/receipts.v1/receiptsv1connect"
	usersv1 "github.com/manzanit0/mcduck/gen/api/users.v1"
	"github.com/manzanit0/mcduck/gen/api/users.v1/usersv1connect"
	"github.com/manzanit0/mcduck/pkg/auth"
)

type ReceiptUploader interface {
	UploadFromChat(ctx context.Context, data []byte, chatID int64) (uint64, error)
}

type APIUploader struct {
	users    usersv1connect.UsersServiceClient
	receipts receiptsv1connect.ReceiptsServiceClient
}

var _ ReceiptUploader = (*APIUploader)(nil)

func NewReceiptUploader(users usersv1connect.UsersServiceClient, receipts receiptsv1connect.ReceiptsServiceClient) *APIUploader {
	return &APIUploader{users: users, receipts: receipts}
}

func (s APIUploader) UploadFromChat(ctx context.Context, data []byte, chatID int64) (uint64, error) {
	getUserReq := connect.Request[usersv1.GetUserRequest]{
		Msg: &usersv1.GetUserRequest{
			TelegramChatId: chatID,
		},
	}

	adminEmail := "admin@mcduck.com"
	token, err := auth.GenerateJWT(adminEmail)
	if err != nil {
		return 0, fmt.Errorf("generate admin JTW: %w", err)
	}

	getUserReq.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := s.users.GetUser(ctx, &getUserReq)
	if err != nil {
		return 0, fmt.Errorf("get user: %w", err)
	}

	createReceiptReq := connect.Request[receiptsv1.CreateReceiptsRequest]{
		Msg: &receiptsv1.CreateReceiptsRequest{
			ReceiptFiles: [][]byte{data},
		},
	}

	onBehalfOf := resp.Msg.User.Email
	token, err = auth.GenerateJWT(onBehalfOf)
	if err != nil {
		return 0, fmt.Errorf("generate user JWT: %w", err)
	}

	createReceiptReq.Header().Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := s.receipts.CreateReceipts(ctx, &createReceiptReq)
	if err != nil {
		return 0, fmt.Errorf("create receipt: %w", err)
	}

	return res.Msg.Receipts[0].GetId(), nil
}
