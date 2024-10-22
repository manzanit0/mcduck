package servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jmoiron/sqlx"
	expensesv1 "github.com/manzanit0/mcduck/gen/api/expenses.v1"
	"github.com/manzanit0/mcduck/gen/api/expenses.v1/expensesv1connect"
	"github.com/manzanit0/mcduck/internal/mcduck"
	"github.com/manzanit0/mcduck/pkg/auth"
	"github.com/manzanit0/mcduck/pkg/xtrace"
)

type expensesServer struct {
	Expenses *mcduck.ExpenseRepository
}

var _ expensesv1connect.ExpensesServiceClient = &expensesServer{}

func NewExpensesServer(db *sqlx.DB) *expensesServer {
	return &expensesServer{
		Expenses: mcduck.NewExpenseRepository(db),
	}
}

// CreateExpense implements expensesv1connect.ExpensesServiceClient.
func (e *expensesServer) CreateExpense(ctx context.Context, req *connect.Request[expensesv1.CreateExpenseRequest]) (*connect.Response[expensesv1.CreateExpenseResponse], error) {
	span := trace.SpanFromContext(ctx)
	email := auth.MustGetUserEmailConnect(ctx)

	expenseID, err := e.Expenses.CreateExpense(ctx, mcduck.CreateExpenseRequest{
		UserEmail:   email,
		Date:        req.Msg.Date.AsTime(),
		Amount:      mcduck.ConvertToDollar(int32(req.Msg.Amount)),
		ReceiptID:   req.Msg.ReceiptId,
		Category:    req.Msg.Category,
		Subcategory: req.Msg.Subcategory,
		Description: req.Msg.Description,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to create expense", "error", err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to create expense: %w", err))
	}

	var category string
	if req.Msg.Category != nil {
		category = *req.Msg.Category
	}

	var subcategory string
	if req.Msg.Subcategory != nil {
		category = *req.Msg.Subcategory
	}

	var description string
	if req.Msg.Description != nil {
		category = *req.Msg.Description
	}

	res := connect.NewResponse(&expensesv1.CreateExpenseResponse{
		Expense: &expensesv1.Expense{
			Id:          uint64(expenseID),
			ReceiptId:   req.Msg.ReceiptId,
			Amount:      req.Msg.Amount,
			Date:        req.Msg.Date,
			Category:    category,
			Subcategory: subcategory,
			Description: description,
		},
	})
	return res, nil
}

// DeleteExpense implements expensesv1connect.ExpensesServiceClient.
func (e *expensesServer) DeleteExpense(ctx context.Context, req *connect.Request[expensesv1.DeleteExpenseRequest]) (*connect.Response[expensesv1.DeleteExpenseResponse], error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("expense.id", int(req.Msg.Id)))

	err := e.Expenses.DeleteExpense(ctx, int64(req.Msg.Id))
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete expense", "error", err.Error())
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to delete expense: %w", err))
	}

	res := connect.NewResponse(&expensesv1.DeleteExpenseResponse{})
	return res, nil
}

// ListExpenses implements expensesv1connect.ExpensesServiceClient.
func (e *expensesServer) ListExpenses(context.Context, *connect.Request[expensesv1.ListExpensesRequest]) (*connect.Response[expensesv1.ListExpensesResponse], error) {
	panic("unimplemented")
}

// UpdateExpense implements expensesv1connect.ExpensesServiceClient.
func (e *expensesServer) UpdateExpense(ctx context.Context, req *connect.Request[expensesv1.UpdateExpenseRequest]) (*connect.Response[expensesv1.UpdateExpenseResponse], error) {
	ctx, span := xtrace.GetSpan(ctx)

	_, err := e.Expenses.FindExpense(ctx, int64(req.Msg.Id))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("expense with id %d doesn't exist", req.Msg.Id))
	} else if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to find expense: %w", err))
	}

	var date *time.Time
	if req.Msg.Date != nil {
		d := req.Msg.Date.AsTime()
		date = &d
	}

	var amount *float32
	if req.Msg.Amount != nil {
		a := mcduck.ConvertToDollar(int32(*req.Msg.Amount))
		amount = &a
	}

	err = e.Expenses.UpdateExpense(ctx, mcduck.UpdateExpenseRequest{
		ID:          int64(req.Msg.Id),
		Date:        date,
		Amount:      amount,
		Category:    req.Msg.Category,
		Subcategory: req.Msg.Subcategory,
		Description: req.Msg.Description,
		ReceiptID:   req.Msg.ReceiptId,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to update expense", "error", err.Error())
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to update expense: %w", err))
	}

	// TODO: merge this with the UpdateExpense SQL call via RETURNING.
	exp, err := e.Expenses.FindExpense(ctx, int64(req.Msg.Id))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("unable to find expense: %w", err))
	}

	var receiptID *uint64
	if exp.ReceiptID != 0 {
		receiptID = &exp.ReceiptID
	}

	res := connect.NewResponse(&expensesv1.UpdateExpenseResponse{
		Expense: &expensesv1.Expense{
			Id:          exp.ID,
			ReceiptId:   receiptID,
			Amount:      uint64(mcduck.ConvertToCents(exp.Amount)),
			Date:        timestamppb.New(exp.Date),
			Category:    exp.Category,
			Subcategory: exp.Subcategory,
			Description: exp.Description,
		},
	})

	return res, nil
}
