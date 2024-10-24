package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"connectrpc.com/otelconnect"
	"github.com/rs/cors"

	"github.com/manzanit0/mcduck/cmd/dots/servers"
	"github.com/manzanit0/mcduck/gen/api/auth.v1/authv1connect"
	"github.com/manzanit0/mcduck/gen/api/expenses.v1/expensesv1connect"
	"github.com/manzanit0/mcduck/gen/api/receipts.v1/receiptsv1connect"
	"github.com/manzanit0/mcduck/gen/api/users.v1/usersv1connect"
	"github.com/manzanit0/mcduck/pkg/auth"
	"github.com/manzanit0/mcduck/pkg/micro"
	"github.com/manzanit0/mcduck/pkg/pubsub"
	"github.com/manzanit0/mcduck/pkg/tgram"
	"github.com/manzanit0/mcduck/pkg/xhttp"
	"github.com/manzanit0/mcduck/pkg/xlog"
	"github.com/manzanit0/mcduck/pkg/xsql"
	"github.com/manzanit0/mcduck/pkg/xtrace"
)

func main() {
	if err := run(); err != nil {
		slog.Error("exiting server", "error", err.Error())
		os.Exit(1)
	}
}

func run() error {
	xlog.InitSlog()

	shutdown, err := xtrace.SetupOTelHTTP(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		err = shutdown(context.Background())
		if err != nil {
			slog.Error("error shutting down OTel SDK")
		}
	}()

	dbx, err := xsql.OpenFromEnv()
	if err != nil {
		return err
	}
	defer xsql.Close(dbx)

	tgramToken := micro.MustGetEnv("TELEGRAM_BOT_TOKEN")
	tgramClient := tgram.NewClient(xhttp.NewClient(), tgramToken)

	natsURL := micro.MustGetEnv("NATS_URL")
	js, _, err := pubsub.NewStream(context.TODO(), natsURL, pubsub.DefaultStreamName, "events.receipts.v1.ReceiptCreated")
	if err != nil {
		return err
	}

	otelInterceptor, err := otelconnect.NewInterceptor(otelconnect.WithTrustRemote())
	if err != nil {
		return err
	}

	authInterceptor := auth.AuthenticationInterceptor()
	traceEnhancer := xtrace.SpanEnhancerInterceptor()

	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "pong"}`))
	}))

	mux.Handle(authv1connect.NewAuthServiceHandler(
		servers.NewAuthServer(dbx, tgramClient),
		connect.WithInterceptors(otelInterceptor, traceEnhancer),
	))

	mux.Handle(receiptsv1connect.NewReceiptsServiceHandler(
		servers.NewReceiptsServer(dbx, tgramClient, js),
		connect.WithInterceptors(otelInterceptor, authInterceptor, traceEnhancer),
	))

	mux.Handle(usersv1connect.NewUsersServiceHandler(
		servers.NewUsersServer(dbx),
		connect.WithInterceptors(otelInterceptor, authInterceptor, traceEnhancer),
	))

	mux.Handle(expensesv1connect.NewExpensesServiceHandler(
		servers.NewExpensesServer(dbx),
		connect.WithInterceptors(otelInterceptor, authInterceptor, traceEnhancer),
	))

	return micro.RunGracefully(withCORS(mux))
}

// withCORS adds CORS support to a Connect HTTP handler.
func withCORS(h http.Handler) http.Handler {
	allowedOrigins := micro.MustGetEnv("ALLOWED_ORIGINS")
	slog.Info("allowed origins: " + allowedOrigins)

	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{allowedOrigins},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}
