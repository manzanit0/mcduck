package servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/jmoiron/sqlx"
	receiptsv1 "github.com/manzanit0/mcduck/gen/api/receipts.v1"
	"github.com/manzanit0/mcduck/gen/api/receipts.v1/receiptsv1connect"
	receiptsevv1 "github.com/manzanit0/mcduck/gen/events/receipts.v1"
	"github.com/manzanit0/mcduck/internal/mcduck"
	"github.com/manzanit0/mcduck/pkg/auth"
	"github.com/manzanit0/mcduck/pkg/pubsub"
	"github.com/manzanit0/mcduck/pkg/tgram"
	"github.com/manzanit0/mcduck/pkg/xtrace"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type receiptsServer struct {
	Telegram tgram.Client
	Receipts *mcduck.ReceiptRepository
	Expenses *mcduck.ExpenseRepository
	js       jetstream.JetStream
}

var _ receiptsv1connect.ReceiptsServiceClient = &receiptsServer{}

func NewReceiptsServer(db *sqlx.DB, t tgram.Client, js jetstream.JetStream) receiptsv1connect.ReceiptsServiceClient {
	return &receiptsServer{
		Telegram: t,
		Receipts: mcduck.NewReceiptRepository(db),
		Expenses: mcduck.NewExpenseRepository(db),
		js:       js,
	}
}

func (s *receiptsServer) CreateReceipts(ctx context.Context, req *connect.Request[receiptsv1.CreateReceiptsRequest]) (*connect.Response[receiptsv1.CreateReceiptsResponse], error) {
	span := trace.SpanFromContext(ctx)

	email := auth.MustGetUserEmailConnect(ctx)

	ch := make(chan *mcduck.Receipt, len(req.Msg.ReceiptFiles))

	g, ctx := errgroup.WithContext(ctx)
	for i, file := range req.Msg.ReceiptFiles {
		g.Go(func() error {
			ctx, span := xtrace.StartSpan(ctx, "Process Receipt")
			defer span.End()

			// TODO: we should do a batch insert to make it an all or nothing.
			created, err := s.Receipts.CreateReceipt(ctx, mcduck.CreateReceiptRequest{
				Image: file,
				Email: email,
			})
			if err != nil {
				slog.ErrorContext(ctx, "failed to insert receipt", "error", err.Error(), "index", i, "email", email)
				span.SetStatus(codes.Error, err.Error())
				return fmt.Errorf("create receipt: %w", err)
			}

			data, topic, err := pubsub.MarshalProto(&receiptsevv1.ReceiptCreated{
				Receipt: &receiptsevv1.Receipt{
					Id:   created.ID,
					File: file,
				},
				UserEmail: email,
			})
			if err != nil {
				slog.ErrorContext(ctx, "failed to marshal receipt to event bytes", "error", err.Error(), "index", i)
				span.SetStatus(codes.Error, err.Error())
				return fmt.Errorf("marshal receipt: %w", err)
			}

			// TODO: this must be done within the Database transaction
			_, publishSpan := xtrace.StartSpan(ctx, fmt.Sprintf("Send %s message", topic))
			slog.InfoContext(ctx, "emitting event", "topic", topic)
			_, err = s.js.Publish(ctx, topic, data)
			if err != nil {
				slog.ErrorContext(ctx, "failed to send ReceiptCreated event", "error", err.Error(), "index", i)
				publishSpan.SetStatus(codes.Error, err.Error())
				publishSpan.End()
				return fmt.Errorf("send event: %w", err)
			}
			publishSpan.End()

			ch <- created

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		slog.ErrorContext(ctx, "create receipt", "error", err.Error())
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	close(ch)

	res := connect.NewResponse(&receiptsv1.CreateReceiptsResponse{})

	for e := range ch {
		res.Msg.Receipts = append(res.Msg.Receipts, &receiptsv1.CreatedReceipt{
			Id:     e.ID,
			Status: receiptsv1.ReceiptStatus_RECEIPT_STATUS_UPLOADED,
		})
	}

	return res, nil
}

func (s *receiptsServer) UpdateReceipt(ctx context.Context, req *connect.Request[receiptsv1.UpdateReceiptRequest]) (*connect.Response[receiptsv1.UpdateReceiptResponse], error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("receipt.id", int(req.Msg.Id)))

	_, err := s.Receipts.GetReceipt(ctx, req.Msg.Id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("receipt with id %d doesn't exist", req.Msg.Id))
	} else if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to find receipt: %w", err))
	}

	var date *time.Time
	if req.Msg.Date != nil {
		d := req.Msg.Date.AsTime()
		date = &d
	}

	dto := mcduck.UpdateReceiptRequest{
		ID:            req.Msg.Id,
		Vendor:        req.Msg.Vendor,
		PendingReview: req.Msg.PendingReview,
		Date:          date,
	}

	err = s.Receipts.UpdateReceipt(ctx, dto)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update receipt", "error", err.Error())
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to update receipt: %w", err))
	}

	res := connect.NewResponse(&receiptsv1.UpdateReceiptResponse{})
	return res, nil
}

func (s *receiptsServer) DeleteReceipt(ctx context.Context, req *connect.Request[receiptsv1.DeleteReceiptRequest]) (*connect.Response[receiptsv1.DeleteReceiptResponse], error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("receipt.id", int(req.Msg.Id)))

	err := s.Receipts.DeleteReceipt(ctx, req.Msg.Id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete receipt", "error", err.Error())
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to delete receipt: %w", err))
	}

	res := connect.NewResponse(&receiptsv1.DeleteReceiptResponse{})
	return res, nil
}

func (s *receiptsServer) ListReceipts(ctx context.Context, req *connect.Request[receiptsv1.ListReceiptsRequest]) (*connect.Response[receiptsv1.ListReceiptsResponse], error) {
	userEmail := auth.MustGetUserEmailConnect(ctx)

	listCtx, span := xtrace.StartSpan(ctx, "List Receipts")
	defer span.End()

	var since mcduck.SinceFilter
	switch req.Msg.Since {
	case receiptsv1.ListReceiptsSince_LIST_RECEIPTS_SINCE_CURRENT_MONTH:
		since = mcduck.SinceFilterCurrentMonth
	case receiptsv1.ListReceiptsSince_LIST_RECEIPTS_SINCE_PREVIOUS_MONTH:
		since = mcduck.SinceFilterPreviousMonth
	default:
	}

	var status mcduck.ReceiptStatus
	switch req.Msg.Status {
	case receiptsv1.ReceiptStatus_RECEIPT_STATUS_UPLOADED:
		status = mcduck.ReceiptStatusUploaded

	case receiptsv1.ReceiptStatus_RECEIPT_STATUS_PENDING_REVIEW:
		status = mcduck.ReceiptStatusUploaded

	case receiptsv1.ReceiptStatus_RECEIPT_STATUS_FAILED_PREPROCESSING:
		status = mcduck.ReceiptStatusFailedPreprocessing

	case receiptsv1.ReceiptStatus_RECEIPT_STATUS_REVIEWED:
		status = mcduck.ReceiptStatusReviewed

	default:
	}

	receipts, err := s.Receipts.ListReceiptsX(listCtx, mcduck.ListFilter{
		Email:  userEmail,
		Since:  since,
		Status: status,
	})
	if err != nil {
		slog.ErrorContext(listCtx, "failed to list receipts", "error", err.Error())
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to list receipts: %w", err))
	}

	span.SetAttributes(attribute.Int("receipts.initial_amount", len(receipts)))

	var resReceipts []*receiptsv1.Receipt
	for _, r := range receipts {
		resReceipts = append(resReceipts, &receiptsv1.Receipt{
			Id:          r.ID,
			Status:      mapReceiptStatus(&r),
			Vendor:      r.Vendor,
			Date:        timestamppb.New(r.Date),
			TotalAmount: r.TotalAmount,
		})
	}

	res := connect.NewResponse(&receiptsv1.ListReceiptsResponse{Receipts: resReceipts})
	return res, nil
}

func (s *receiptsServer) GetReceipt(ctx context.Context, req *connect.Request[receiptsv1.GetReceiptRequest]) (*connect.Response[receiptsv1.GetReceiptResponse], error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("receipt.id", int(req.Msg.Id)))

	receipt, err := s.Receipts.GetReceipt(ctx, req.Msg.Id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("receipt with id %d doesn't exist", req.Msg.Id))
	} else if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to get receipt", "error", err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get receipt: %w", err))
	}

	expenses, err := s.Expenses.ListExpensesForReceipt(ctx, req.Msg.Id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to list expenses for receipt", "error", err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to get expenses for receipt: %w", err))
	}

	res := connect.NewResponse(&receiptsv1.GetReceiptResponse{
		Receipt: &receiptsv1.FullReceipt{
			Id:       receipt.ID,
			Status:   mapReceiptStatus(receipt),
			Vendor:   receipt.Vendor,
			Date:     timestamppb.New(receipt.Date),
			File:     receipt.Image,
			Expenses: mapExpenses(expenses),
		},
	})

	return res, nil
}

func mapExpenses(expenses []mcduck.Expense) []*receiptsv1.Expense {
	resExpenses := make([]*receiptsv1.Expense, len(expenses))
	for i, e := range expenses {
		resExp := receiptsv1.Expense{
			Id:          e.ID,
			Date:        timestamppb.New(e.Date),
			Category:    e.Category,
			Subcategory: e.Subcategory,
			Description: e.Description,
			Amount:      e.Amount,
		}

		resExpenses[i] = &resExp
	}

	return resExpenses
}

func mapReceiptStatus(r *mcduck.Receipt) receiptsv1.ReceiptStatus {
	switch r.Status {
	case mcduck.ReceiptStatusUploaded:
		return receiptsv1.ReceiptStatus_RECEIPT_STATUS_UPLOADED

	case mcduck.ReceiptStatusFailedPreprocessing:
		return receiptsv1.ReceiptStatus_RECEIPT_STATUS_FAILED_PREPROCESSING

	case mcduck.ReceiptStatusPendingReview:
		return receiptsv1.ReceiptStatus_RECEIPT_STATUS_PENDING_REVIEW

	case mcduck.ReceiptStatusReviewed:
		return receiptsv1.ReceiptStatus_RECEIPT_STATUS_REVIEWED
	}

	return receiptsv1.ReceiptStatus_RECEIPT_STATUS_UNSPECIFIED
}
