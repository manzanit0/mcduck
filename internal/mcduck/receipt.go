package mcduck

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/manzanit0/mcduck/internal/parser"
	"github.com/manzanit0/mcduck/pkg/xsql"
	"github.com/manzanit0/mcduck/pkg/xtrace"

	receiptsevv1 "github.com/manzanit0/mcduck/gen/events/receipts.v1"
)

type Receipt struct {
	ID            uint64
	PendingReview bool
	Status        string
	Image         []byte
	Vendor        string
	UserEmail     string
	Date          time.Time
	TotalAmount   uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type dbReceipt struct {
	ID            uint64  `db:"id"`
	PendingReview bool    `db:"pending_review"`
	Status        string  `db:"status"`
	Image         []byte  `db:"receipt_image"`
	UserEmail     string  `db:"user_email"`
	Vendor        *string `db:"vendor"`
	TotalAmount   *uint64 `db:"total_amount"`

	Date      time.Time `db:"receipt_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r *dbReceipt) MapReceipt() *Receipt {
	receipt := Receipt{
		ID:            r.ID,
		PendingReview: r.PendingReview,
		Image:         r.Image,
		Date:          r.Date,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		Status:        r.Status,
		UserEmail:     r.UserEmail,
	}

	if r.Vendor != nil {
		receipt.Vendor = *r.Vendor
	}

	if r.TotalAmount != nil {
		receipt.TotalAmount = *r.TotalAmount
	}

	return &receipt
}

type ReceiptRepository struct {
	dbx *sqlx.DB
}

func NewReceiptRepository(dbx *sqlx.DB) *ReceiptRepository {
	return &ReceiptRepository{dbx: dbx}
}

type CreateReceiptRequest struct {
	Amount      uint64
	Description string
	Vendor      string
	Image       []byte
	Date        time.Time
	Email       string
}

func (r *ReceiptRepository) CreateReceipt(ctx context.Context, input CreateReceiptRequest) (*Receipt, error) {
	ctx, span := xtrace.StartSpan(ctx, "Create Receipt")
	defer span.End()

	if len(input.Image) == 0 {
		return nil, fmt.Errorf("empty receipt")
	}

	txn, err := r.dbx.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	defer xsql.TxClose(txn)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.
		Insert("receipts").
		Columns("receipt_image", "pending_review", "user_email", "receipt_date", "vendor").
		Values(input.Image, true, input.Email, input.Date, input.Vendor).
		Suffix(`RETURNING id, pending_review, receipt_date, vendor, user_email`)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build query: %w", err)
	}

	var record dbReceipt
	err = txn.GetContext(ctx, &record, query, args...)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}

	if input.Amount > 0 {
		e := ExpensesBatch{
			UserEmail: input.Email,
			Records: []Expense{{
				ReceiptID:   record.ID,
				Date:        input.Date,
				Amount:      input.Amount,
				UserEmail:   input.Email,
				Description: input.Description,
				Category:    "Receipt Upload",
			}},
		}

		err = CreateExpenses(ctx, txn, e)
		if err != nil {
			return nil, fmt.Errorf("unable to insert expenses: %w", err)
		}
	}

	err = txn.Commit()
	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return record.MapReceipt(), nil
}

type UpdateReceiptRequest struct {
	ID            uint64
	Vendor        *string
	PendingReview *bool
	Date          *time.Time
}

func (r *ReceiptRepository) UpdateReceipt(ctx context.Context, e UpdateReceiptRequest) error {
	ctx, span := xtrace.StartSpan(ctx, "Update Receipt")
	defer span.End()

	var shouldUpdate bool
	var shouldUpdateExpenseDates bool

	txn, err := r.dbx.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer xsql.TxClose(txn)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Update("receipts").Where(sq.Eq{"id": e.ID})

	if e.Vendor != nil {
		builder = builder.Set("vendor", *e.Vendor)
		shouldUpdate = true
	}

	if e.PendingReview != nil {
		builder = builder.Set("pending_review", *e.PendingReview)
		if *e.PendingReview {
			builder = builder.Set("status", "pending_review")
		} else {
			builder = builder.Set("status", "reviewed")
		}
		shouldUpdate = true
	}

	if e.Date != nil {
		builder = builder.Set("receipt_date", *e.Date)
		shouldUpdate = true
		shouldUpdateExpenseDates = true
	}

	if !shouldUpdate {
		return nil
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("compile receipts query: %w", err)
	}

	_, err = txn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	if shouldUpdateExpenseDates {
		query, args, err = psql.Update("expenses").Where(sq.Eq{"receipt_id": e.ID}).Set("expense_date", *e.Date).ToSql()
		if err != nil {
			return fmt.Errorf("compile expenses query: %w", err)
		}

		_, err = txn.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("execute expenses query: %w", err)
		}
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (r *ExpenseRepository) UpdateReceiptWithTxn(ctx context.Context, txn *sqlx.Tx, e UpdateReceiptRequest) error {
	ctx, span := xtrace.StartSpan(ctx, "Update Receipt")
	defer span.End()

	var shouldUpdate bool
	var shouldUpdateExpenseDates bool

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Update("receipts").Where(sq.Eq{"id": e.ID})

	if e.Vendor != nil {
		builder = builder.Set("vendor", *e.Vendor)
		shouldUpdate = true
	}

	if e.PendingReview != nil {
		builder = builder.Set("pending_review", *e.PendingReview)
		if *e.PendingReview {
			builder = builder.Set("status", "pending_review")
		} else {
			builder = builder.Set("status", "reviewed")
		}
		shouldUpdate = true
	}

	if e.Date != nil {
		builder = builder.Set("receipt_date", *e.Date)
		shouldUpdate = true
		shouldUpdateExpenseDates = true
	}

	if !shouldUpdate {
		return nil
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("compile receipts query: %w", err)
	}

	_, err = txn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	if shouldUpdateExpenseDates {
		query, args, err = psql.Update("expenses").Where(sq.Eq{"receipt_id": e.ID}).Set("expense_date", *e.Date).ToSql()
		if err != nil {
			return fmt.Errorf("compile expenses query: %w", err)
		}

		_, err = txn.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("execute expenses query: %w", err)
		}
	}

	return nil
}

func (r *ReceiptRepository) MarkFailedToProcess(ctx context.Context, receiptID uint64) error {
	ctx, span := xtrace.StartSpan(ctx, "Update Receipt")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.
		Update("receipts").
		Set("status", "failed_preprocessing").
		Where(sq.Eq{"id": receiptID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("compile receipts query: %w", err)
	}

	_, err = r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query: %w", err)
	}

	return nil
}

func (r *ReceiptRepository) ListReceipts(ctx context.Context, email string) ([]Receipt, error) {
	ctx, span := xtrace.StartSpan(ctx, "List Receipts")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select("id", "vendor", "pending_review", "receipt_date", "status").
		From("receipts").
		Where(sq.Eq{"user_email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("compile query: %w", err)
	}

	var receipts []dbReceipt
	err = r.dbx.SelectContext(ctx, &receipts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select receipts: %w", err)
	}

	var domainReceipts []Receipt
	for _, receipt := range receipts {
		domainReceipts = append(domainReceipts, *receipt.MapReceipt())
	}

	return domainReceipts, nil
}

func (r *ReceiptRepository) ListReceiptsCurrentMonth(ctx context.Context, email string) ([]Receipt, error) {
	ctx, span := xtrace.StartSpan(ctx, "List Receipts for Current Month")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select("id", "vendor", "pending_review", "receipt_date", "status").
		From("receipts").
		Where(sq.And{
			sq.Eq{"user_email": email},
			sq.Expr("receipt_date >= date_trunc('month',current_date)"),
			sq.Expr("receipt_date < date_trunc('month',current_date) + INTERVAL '1' MONTH"),
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("compile query: %w", err)
	}

	var receipts []dbReceipt
	err = r.dbx.SelectContext(ctx, &receipts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select receipts: %w", err)
	}

	var domainReceipts []Receipt
	for _, receipt := range receipts {
		domainReceipts = append(domainReceipts, *receipt.MapReceipt())
	}

	return domainReceipts, nil
}

func (r *ReceiptRepository) ListReceiptsPreviousMonth(ctx context.Context, email string) ([]Receipt, error) {
	ctx, span := xtrace.StartSpan(ctx, "List Receipts for Previous Month")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select("id", "vendor", "pending_review", "receipt_date", "status").
		From("receipts").
		Where(sq.And{
			sq.Eq{"user_email": email},
			sq.Expr("receipt_date >= date_trunc('month',current_date) - INTERVAL '1' MONTH"),
			sq.Expr("receipt_date < date_trunc('month',current_date)"),
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("compile query: %w", err)
	}

	var receipts []dbReceipt
	err = r.dbx.SelectContext(ctx, &receipts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select receipts: %w", err)
	}

	var domainReceipts []Receipt
	for _, receipt := range receipts {
		domainReceipts = append(domainReceipts, *receipt.MapReceipt())
	}

	return domainReceipts, nil
}

func (r *ReceiptRepository) ListReceiptsPendingReview(ctx context.Context, email string) ([]Receipt, error) {
	ctx, span := xtrace.StartSpan(ctx, "List Receipts Pending Review")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select("id", "vendor", "pending_review", "receipt_date", "status").
		From("receipts").
		Where(sq.And{
			sq.Eq{"user_email": email},
			sq.Eq{"pending_review": true},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("compile query: %w", err)
	}

	var receipts []dbReceipt
	err = r.dbx.SelectContext(ctx, &receipts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select receipts: %w", err)
	}

	var domainReceipts []Receipt
	for _, receipt := range receipts {
		domainReceipts = append(domainReceipts, *receipt.MapReceipt())
	}

	return domainReceipts, nil
}

func (r *ReceiptRepository) GetReceipt(ctx context.Context, receiptID uint64) (*Receipt, error) {
	ctx, span := xtrace.StartSpan(ctx, "Get Single Receipt")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select("id", "vendor", "pending_review", "created_at", "receipt_image", "user_email", "receipt_date", "status").
		From("receipts").
		Where(sq.Eq{"id": receiptID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("compile query: %w", err)
	}

	var receipt dbReceipt
	err = r.dbx.GetContext(ctx, &receipt, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select receipts: %w", err)
	}

	return receipt.MapReceipt(), nil
}

func (r *ReceiptRepository) GetReceiptImage(ctx context.Context, receiptID uint64) ([]byte, error) {
	ctx, span := xtrace.StartSpan(ctx, "Get Receipt Image")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select("receipt_image").
		From("receipts").
		Where(sq.Eq{"id": receiptID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("compile query: %w", err)
	}

	var receipt dbReceipt
	err = r.dbx.GetContext(ctx, &receipt, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select receipts: %w", err)
	}

	return receipt.Image, nil
}

func (r *ReceiptRepository) DeleteReceipt(ctx context.Context, id uint64) error {
	ctx, span := xtrace.StartSpan(ctx, "Delete Receipt")
	defer span.End()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	txn, err := r.dbx.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer xsql.TxClose(txn)

	query, args, err := psql.Delete("receipts").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("unable to build query: %w", err)
	}

	_, err = r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to execute query to delete receipt: %w", err)
	}

	query, args, err = psql.Delete("expenses").Where(sq.Eq{"receipt_id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("unable to build query: %w", err)
	}

	_, err = r.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("unable to execute query to delete expenses: %w", err)
	}

	err = txn.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

type ReceiptAugmentor interface {
	AugmentCreatedReceipt(ctx context.Context, ev *receiptsevv1.ReceiptCreated) error
}

type AIaugmentor struct {
	DB     *sqlx.DB
	Parser parser.ReceiptParser
}

func (a *AIaugmentor) AugmentCreatedReceipt(ctx context.Context, ev *receiptsevv1.ReceiptCreated) error {
	data := ev.Receipt.File
	parsedReceipt, _, err := a.Parser.ExtractReceipt(ctx, data)
	if err != nil {
		return fmt.Errorf("parser receipt: %w", err)
	}

	parsedTime, err := time.Parse("02/01/2006", parsedReceipt.PurchaseDate)
	if err != nil {
		slog.Info("failed to parse receipt date. Defaulting to 'now' ", "error", err.Error())
		parsedTime = time.Now()
	}

	txn, err := a.DB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer xsql.TxClose(txn)

	receiptRepo := NewExpenseRepository(a.DB)

	pendingReview := true
	err = receiptRepo.UpdateReceiptWithTxn(ctx, txn, UpdateReceiptRequest{
		ID:            ev.Receipt.Id,
		Vendor:        &parsedReceipt.Vendor,
		Date:          &parsedTime,
		PendingReview: &pendingReview,
	})
	if err != nil {
		return fmt.Errorf("update receipt: %w", err)
	}

	err = CreateExpenses(ctx, txn, ExpensesBatch{
		Records: []Expense{
			{
				Date:        parsedTime,
				Amount:      ConvertToCents(float32(parsedReceipt.Amount)),
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
		return fmt.Errorf("create expenses: %w", err)
	}

	if err = txn.Commit(); err != nil {
		return fmt.Errorf("commit txn: %w", err)
	}

	return nil
}

type SinceFilter string

const (
	SinceFilterCurrentMonth  SinceFilter = "current_month"
	SinceFilterPreviousMonth SinceFilter = "previous_month"
)

type ReceiptStatus string

const (
	ReceiptStatusUploaded            = "uploaded"
	ReceiptStatusPendingReview       = "pending_review"
	ReceiptStatusFailedPreprocessing = "failed_preprocessing"
	ReceiptStatusReviewed            = "reviewed"
)

type ListFilter struct {
	Email  string
	Since  SinceFilter
	Status ReceiptStatus
}

func (r *ReceiptRepository) ListReceiptsX(ctx context.Context, filter ListFilter) ([]Receipt, error) {
	ctx, span := xtrace.StartSpan(ctx, "List Receipts for Current Month")
	defer span.End()

	andFilters := sq.And{sq.Eq{"r.user_email": filter.Email}}

	switch filter.Since {
	case SinceFilterCurrentMonth:
		andFilters = append(andFilters,
			sq.Expr("receipt_date >= date_trunc('month',current_date)"),
			sq.Expr("receipt_date < date_trunc('month',current_date) + INTERVAL '1' MONTH"),
		)
	case SinceFilterPreviousMonth:
		andFilters = append(andFilters,
			sq.Expr("receipt_date >= date_trunc('month',current_date) - INTERVAL '1' MONTH"),
			sq.Expr("receipt_date < date_trunc('month',current_date)"),
		)
	}

	if filter.Status != "" {
		andFilters = append(andFilters, sq.Expr("receipt_status = ?", filter.Status))
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Select(
			"r.id",
			"r.pending_review",
			"r.status",
			"r.user_email",
			"r.vendor",
			"r.receipt_date",
			"(SELECT SUM(amount) AS total_amount FROM expenses WHERE expenses.receipt_id=r.id)",
		).
		From("receipts r").
		Where(andFilters).
		OrderBy("r.receipt_date DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("compile query: %w", err)
	}

	var receipts []dbReceipt
	err = r.dbx.SelectContext(ctx, &receipts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select receipts: %w", err)
	}

	var domainReceipts []Receipt
	for _, receipt := range receipts {
		domainReceipts = append(domainReceipts, *receipt.MapReceipt())
	}

	return domainReceipts, nil
}
