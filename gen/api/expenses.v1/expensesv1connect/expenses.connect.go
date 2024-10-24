// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/expenses.v1/expenses.proto

package expensesv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	expenses_v1 "github.com/manzanit0/mcduck/gen/api/expenses.v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// ExpensesServiceName is the fully-qualified name of the ExpensesService service.
	ExpensesServiceName = "api.expenses.v1.ExpensesService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ExpensesServiceCreateExpenseProcedure is the fully-qualified name of the ExpensesService's
	// CreateExpense RPC.
	ExpensesServiceCreateExpenseProcedure = "/api.expenses.v1.ExpensesService/CreateExpense"
	// ExpensesServiceUpdateExpenseProcedure is the fully-qualified name of the ExpensesService's
	// UpdateExpense RPC.
	ExpensesServiceUpdateExpenseProcedure = "/api.expenses.v1.ExpensesService/UpdateExpense"
	// ExpensesServiceDeleteExpenseProcedure is the fully-qualified name of the ExpensesService's
	// DeleteExpense RPC.
	ExpensesServiceDeleteExpenseProcedure = "/api.expenses.v1.ExpensesService/DeleteExpense"
	// ExpensesServiceListExpensesProcedure is the fully-qualified name of the ExpensesService's
	// ListExpenses RPC.
	ExpensesServiceListExpensesProcedure = "/api.expenses.v1.ExpensesService/ListExpenses"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	expensesServiceServiceDescriptor             = expenses_v1.File_api_expenses_v1_expenses_proto.Services().ByName("ExpensesService")
	expensesServiceCreateExpenseMethodDescriptor = expensesServiceServiceDescriptor.Methods().ByName("CreateExpense")
	expensesServiceUpdateExpenseMethodDescriptor = expensesServiceServiceDescriptor.Methods().ByName("UpdateExpense")
	expensesServiceDeleteExpenseMethodDescriptor = expensesServiceServiceDescriptor.Methods().ByName("DeleteExpense")
	expensesServiceListExpensesMethodDescriptor  = expensesServiceServiceDescriptor.Methods().ByName("ListExpenses")
)

// ExpensesServiceClient is a client for the api.expenses.v1.ExpensesService service.
type ExpensesServiceClient interface {
	CreateExpense(context.Context, *connect.Request[expenses_v1.CreateExpenseRequest]) (*connect.Response[expenses_v1.CreateExpenseResponse], error)
	UpdateExpense(context.Context, *connect.Request[expenses_v1.UpdateExpenseRequest]) (*connect.Response[expenses_v1.UpdateExpenseResponse], error)
	DeleteExpense(context.Context, *connect.Request[expenses_v1.DeleteExpenseRequest]) (*connect.Response[expenses_v1.DeleteExpenseResponse], error)
	ListExpenses(context.Context, *connect.Request[expenses_v1.ListExpensesRequest]) (*connect.Response[expenses_v1.ListExpensesResponse], error)
}

// NewExpensesServiceClient constructs a client for the api.expenses.v1.ExpensesService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewExpensesServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ExpensesServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &expensesServiceClient{
		createExpense: connect.NewClient[expenses_v1.CreateExpenseRequest, expenses_v1.CreateExpenseResponse](
			httpClient,
			baseURL+ExpensesServiceCreateExpenseProcedure,
			connect.WithSchema(expensesServiceCreateExpenseMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateExpense: connect.NewClient[expenses_v1.UpdateExpenseRequest, expenses_v1.UpdateExpenseResponse](
			httpClient,
			baseURL+ExpensesServiceUpdateExpenseProcedure,
			connect.WithSchema(expensesServiceUpdateExpenseMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deleteExpense: connect.NewClient[expenses_v1.DeleteExpenseRequest, expenses_v1.DeleteExpenseResponse](
			httpClient,
			baseURL+ExpensesServiceDeleteExpenseProcedure,
			connect.WithSchema(expensesServiceDeleteExpenseMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		listExpenses: connect.NewClient[expenses_v1.ListExpensesRequest, expenses_v1.ListExpensesResponse](
			httpClient,
			baseURL+ExpensesServiceListExpensesProcedure,
			connect.WithSchema(expensesServiceListExpensesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// expensesServiceClient implements ExpensesServiceClient.
type expensesServiceClient struct {
	createExpense *connect.Client[expenses_v1.CreateExpenseRequest, expenses_v1.CreateExpenseResponse]
	updateExpense *connect.Client[expenses_v1.UpdateExpenseRequest, expenses_v1.UpdateExpenseResponse]
	deleteExpense *connect.Client[expenses_v1.DeleteExpenseRequest, expenses_v1.DeleteExpenseResponse]
	listExpenses  *connect.Client[expenses_v1.ListExpensesRequest, expenses_v1.ListExpensesResponse]
}

// CreateExpense calls api.expenses.v1.ExpensesService.CreateExpense.
func (c *expensesServiceClient) CreateExpense(ctx context.Context, req *connect.Request[expenses_v1.CreateExpenseRequest]) (*connect.Response[expenses_v1.CreateExpenseResponse], error) {
	return c.createExpense.CallUnary(ctx, req)
}

// UpdateExpense calls api.expenses.v1.ExpensesService.UpdateExpense.
func (c *expensesServiceClient) UpdateExpense(ctx context.Context, req *connect.Request[expenses_v1.UpdateExpenseRequest]) (*connect.Response[expenses_v1.UpdateExpenseResponse], error) {
	return c.updateExpense.CallUnary(ctx, req)
}

// DeleteExpense calls api.expenses.v1.ExpensesService.DeleteExpense.
func (c *expensesServiceClient) DeleteExpense(ctx context.Context, req *connect.Request[expenses_v1.DeleteExpenseRequest]) (*connect.Response[expenses_v1.DeleteExpenseResponse], error) {
	return c.deleteExpense.CallUnary(ctx, req)
}

// ListExpenses calls api.expenses.v1.ExpensesService.ListExpenses.
func (c *expensesServiceClient) ListExpenses(ctx context.Context, req *connect.Request[expenses_v1.ListExpensesRequest]) (*connect.Response[expenses_v1.ListExpensesResponse], error) {
	return c.listExpenses.CallUnary(ctx, req)
}

// ExpensesServiceHandler is an implementation of the api.expenses.v1.ExpensesService service.
type ExpensesServiceHandler interface {
	CreateExpense(context.Context, *connect.Request[expenses_v1.CreateExpenseRequest]) (*connect.Response[expenses_v1.CreateExpenseResponse], error)
	UpdateExpense(context.Context, *connect.Request[expenses_v1.UpdateExpenseRequest]) (*connect.Response[expenses_v1.UpdateExpenseResponse], error)
	DeleteExpense(context.Context, *connect.Request[expenses_v1.DeleteExpenseRequest]) (*connect.Response[expenses_v1.DeleteExpenseResponse], error)
	ListExpenses(context.Context, *connect.Request[expenses_v1.ListExpensesRequest]) (*connect.Response[expenses_v1.ListExpensesResponse], error)
}

// NewExpensesServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewExpensesServiceHandler(svc ExpensesServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	expensesServiceCreateExpenseHandler := connect.NewUnaryHandler(
		ExpensesServiceCreateExpenseProcedure,
		svc.CreateExpense,
		connect.WithSchema(expensesServiceCreateExpenseMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	expensesServiceUpdateExpenseHandler := connect.NewUnaryHandler(
		ExpensesServiceUpdateExpenseProcedure,
		svc.UpdateExpense,
		connect.WithSchema(expensesServiceUpdateExpenseMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	expensesServiceDeleteExpenseHandler := connect.NewUnaryHandler(
		ExpensesServiceDeleteExpenseProcedure,
		svc.DeleteExpense,
		connect.WithSchema(expensesServiceDeleteExpenseMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	expensesServiceListExpensesHandler := connect.NewUnaryHandler(
		ExpensesServiceListExpensesProcedure,
		svc.ListExpenses,
		connect.WithSchema(expensesServiceListExpensesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/api.expenses.v1.ExpensesService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ExpensesServiceCreateExpenseProcedure:
			expensesServiceCreateExpenseHandler.ServeHTTP(w, r)
		case ExpensesServiceUpdateExpenseProcedure:
			expensesServiceUpdateExpenseHandler.ServeHTTP(w, r)
		case ExpensesServiceDeleteExpenseProcedure:
			expensesServiceDeleteExpenseHandler.ServeHTTP(w, r)
		case ExpensesServiceListExpensesProcedure:
			expensesServiceListExpensesHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedExpensesServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedExpensesServiceHandler struct{}

func (UnimplementedExpensesServiceHandler) CreateExpense(context.Context, *connect.Request[expenses_v1.CreateExpenseRequest]) (*connect.Response[expenses_v1.CreateExpenseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.expenses.v1.ExpensesService.CreateExpense is not implemented"))
}

func (UnimplementedExpensesServiceHandler) UpdateExpense(context.Context, *connect.Request[expenses_v1.UpdateExpenseRequest]) (*connect.Response[expenses_v1.UpdateExpenseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.expenses.v1.ExpensesService.UpdateExpense is not implemented"))
}

func (UnimplementedExpensesServiceHandler) DeleteExpense(context.Context, *connect.Request[expenses_v1.DeleteExpenseRequest]) (*connect.Response[expenses_v1.DeleteExpenseResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.expenses.v1.ExpensesService.DeleteExpense is not implemented"))
}

func (UnimplementedExpensesServiceHandler) ListExpenses(context.Context, *connect.Request[expenses_v1.ListExpensesRequest]) (*connect.Response[expenses_v1.ListExpensesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.expenses.v1.ExpensesService.ListExpenses is not implemented"))
}
