package main

import (
	"context"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/otelconnect"
	"github.com/rs/cors"

	"github.com/manzanit0/mcduck/api/auth.v1/authv1connect"
	"github.com/manzanit0/mcduck/api/receipts.v1/receiptsv1connect"
	"github.com/manzanit0/mcduck/cmd/dots/servers"
	"github.com/manzanit0/mcduck/pkg/auth"
	"github.com/manzanit0/mcduck/pkg/micro"
	"github.com/manzanit0/mcduck/pkg/tgram"
	"github.com/manzanit0/mcduck/pkg/xhttp"
	"github.com/manzanit0/mcduck/pkg/xlog"
	"github.com/manzanit0/mcduck/pkg/xsql"
	"github.com/manzanit0/mcduck/pkg/xtrace"
)

const serviceName = "dots"

func main() {
	xlog.InitSlog()

	tp, err := xtrace.TracerFromEnv(context.Background(), serviceName)
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(context.Background())

	dbx, err := xsql.Open(serviceName)
	if err != nil {
		panic(err)
	}
	defer xsql.Close(dbx.GetSQLX())

	tgramToken := micro.MustGetEnv("TELEGRAM_BOT_TOKEN")
	tgramClient := tgram.NewClient(xhttp.NewClient(), tgramToken)

	otelInterceptor, err := otelconnect.NewInterceptor(otelconnect.WithTrustRemote(), otelconnect.WithoutMetrics())
	if err != nil {
		panic(err)
	}

	authInterceptor := auth.AuthenticationInterceptor()
	traceEnhancer := xtrace.SpanEnhancerInterceptor()

	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "pong"}`))
	}))

	mux.Handle(authv1connect.NewAuthServiceHandler(
		servers.NewAuthServer(dbx.GetSQLX(), tgramClient),
		connect.WithInterceptors(otelInterceptor, traceEnhancer),
	))

	mux.Handle(receiptsv1connect.NewReceiptsServiceHandler(
		servers.NewReceiptsServer(dbx.GetSQLX(), tgramClient),
		connect.WithInterceptors(otelInterceptor, authInterceptor, traceEnhancer),
	))

	if err := micro.RunGracefully(withCORS(mux)); err != nil {
		os.Exit(1)
	}
}

// withCORS adds CORS support to a Connect HTTP handler.
func withCORS(h http.Handler) http.Handler {
	allowedOrigins := micro.MustGetEnv("ALLOWED_ORIGINS")

	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigins},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}