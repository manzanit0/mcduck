// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: api/auth.v1/auth.proto

package authv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	auth_v1 "github.com/manzanit0/mcduck/gen/api/auth.v1"
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
	// AuthServiceName is the fully-qualified name of the AuthService service.
	AuthServiceName = "api.auth.v1.AuthService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// AuthServiceRegisterProcedure is the fully-qualified name of the AuthService's Register RPC.
	AuthServiceRegisterProcedure = "/api.auth.v1.AuthService/Register"
	// AuthServiceLoginProcedure is the fully-qualified name of the AuthService's Login RPC.
	AuthServiceLoginProcedure = "/api.auth.v1.AuthService/Login"
	// AuthServiceConnectTelegramProcedure is the fully-qualified name of the AuthService's
	// ConnectTelegram RPC.
	AuthServiceConnectTelegramProcedure = "/api.auth.v1.AuthService/ConnectTelegram"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	authServiceServiceDescriptor               = auth_v1.File_api_auth_v1_auth_proto.Services().ByName("AuthService")
	authServiceRegisterMethodDescriptor        = authServiceServiceDescriptor.Methods().ByName("Register")
	authServiceLoginMethodDescriptor           = authServiceServiceDescriptor.Methods().ByName("Login")
	authServiceConnectTelegramMethodDescriptor = authServiceServiceDescriptor.Methods().ByName("ConnectTelegram")
)

// AuthServiceClient is a client for the api.auth.v1.AuthService service.
type AuthServiceClient interface {
	Register(context.Context, *connect.Request[auth_v1.RegisterRequest]) (*connect.Response[auth_v1.RegisterResponse], error)
	Login(context.Context, *connect.Request[auth_v1.LoginRequest]) (*connect.Response[auth_v1.LoginResponse], error)
	ConnectTelegram(context.Context, *connect.Request[auth_v1.ConnectTelegramRequest]) (*connect.Response[auth_v1.ConnectTelegramResponse], error)
}

// NewAuthServiceClient constructs a client for the api.auth.v1.AuthService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAuthServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) AuthServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &authServiceClient{
		register: connect.NewClient[auth_v1.RegisterRequest, auth_v1.RegisterResponse](
			httpClient,
			baseURL+AuthServiceRegisterProcedure,
			connect.WithSchema(authServiceRegisterMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		login: connect.NewClient[auth_v1.LoginRequest, auth_v1.LoginResponse](
			httpClient,
			baseURL+AuthServiceLoginProcedure,
			connect.WithSchema(authServiceLoginMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		connectTelegram: connect.NewClient[auth_v1.ConnectTelegramRequest, auth_v1.ConnectTelegramResponse](
			httpClient,
			baseURL+AuthServiceConnectTelegramProcedure,
			connect.WithSchema(authServiceConnectTelegramMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// authServiceClient implements AuthServiceClient.
type authServiceClient struct {
	register        *connect.Client[auth_v1.RegisterRequest, auth_v1.RegisterResponse]
	login           *connect.Client[auth_v1.LoginRequest, auth_v1.LoginResponse]
	connectTelegram *connect.Client[auth_v1.ConnectTelegramRequest, auth_v1.ConnectTelegramResponse]
}

// Register calls api.auth.v1.AuthService.Register.
func (c *authServiceClient) Register(ctx context.Context, req *connect.Request[auth_v1.RegisterRequest]) (*connect.Response[auth_v1.RegisterResponse], error) {
	return c.register.CallUnary(ctx, req)
}

// Login calls api.auth.v1.AuthService.Login.
func (c *authServiceClient) Login(ctx context.Context, req *connect.Request[auth_v1.LoginRequest]) (*connect.Response[auth_v1.LoginResponse], error) {
	return c.login.CallUnary(ctx, req)
}

// ConnectTelegram calls api.auth.v1.AuthService.ConnectTelegram.
func (c *authServiceClient) ConnectTelegram(ctx context.Context, req *connect.Request[auth_v1.ConnectTelegramRequest]) (*connect.Response[auth_v1.ConnectTelegramResponse], error) {
	return c.connectTelegram.CallUnary(ctx, req)
}

// AuthServiceHandler is an implementation of the api.auth.v1.AuthService service.
type AuthServiceHandler interface {
	Register(context.Context, *connect.Request[auth_v1.RegisterRequest]) (*connect.Response[auth_v1.RegisterResponse], error)
	Login(context.Context, *connect.Request[auth_v1.LoginRequest]) (*connect.Response[auth_v1.LoginResponse], error)
	ConnectTelegram(context.Context, *connect.Request[auth_v1.ConnectTelegramRequest]) (*connect.Response[auth_v1.ConnectTelegramResponse], error)
}

// NewAuthServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAuthServiceHandler(svc AuthServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	authServiceRegisterHandler := connect.NewUnaryHandler(
		AuthServiceRegisterProcedure,
		svc.Register,
		connect.WithSchema(authServiceRegisterMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	authServiceLoginHandler := connect.NewUnaryHandler(
		AuthServiceLoginProcedure,
		svc.Login,
		connect.WithSchema(authServiceLoginMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	authServiceConnectTelegramHandler := connect.NewUnaryHandler(
		AuthServiceConnectTelegramProcedure,
		svc.ConnectTelegram,
		connect.WithSchema(authServiceConnectTelegramMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/api.auth.v1.AuthService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case AuthServiceRegisterProcedure:
			authServiceRegisterHandler.ServeHTTP(w, r)
		case AuthServiceLoginProcedure:
			authServiceLoginHandler.ServeHTTP(w, r)
		case AuthServiceConnectTelegramProcedure:
			authServiceConnectTelegramHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedAuthServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAuthServiceHandler struct{}

func (UnimplementedAuthServiceHandler) Register(context.Context, *connect.Request[auth_v1.RegisterRequest]) (*connect.Response[auth_v1.RegisterResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.auth.v1.AuthService.Register is not implemented"))
}

func (UnimplementedAuthServiceHandler) Login(context.Context, *connect.Request[auth_v1.LoginRequest]) (*connect.Response[auth_v1.LoginResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.auth.v1.AuthService.Login is not implemented"))
}

func (UnimplementedAuthServiceHandler) ConnectTelegram(context.Context, *connect.Request[auth_v1.ConnectTelegramRequest]) (*connect.Response[auth_v1.ConnectTelegramResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("api.auth.v1.AuthService.ConnectTelegram is not implemented"))
}