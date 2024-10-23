package micro

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manzanit0/mcduck/pkg/xlog"
	"github.com/manzanit0/mcduck/pkg/xtrace"
)

type ShutdownFunc func(context.Context) error

type Service struct {
	Name          string
	Engine        *gin.Engine
	shutdownFuncs []ShutdownFunc
}

func NewGinService(name string) (Service, error) {
	xlog.InitSlog()

	shutdownOTel, err := xtrace.SetupOTelHTTP(context.Background())
	if err != nil {
		return Service{}, fmt.Errorf("get tracer from env %w", err)
	}

	r := gin.Default()
	r.Use(xlog.EnhanceContext)
	r.Use(xtrace.GinTraceRequests(name))
	r.Use(xtrace.GinEnhanceTraceAttributes())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return Service{Name: name, Engine: r, shutdownFuncs: []ShutdownFunc{shutdownOTel}}, nil
}

func (s *Service) Run() error {
	for _, fn := range s.shutdownFuncs {
		defer func() {
			err := fn(context.Background())
			if err != nil {
				slog.Error("error shutting down thing", "error", err.Error())
			}
		}()
	}

	return RunGracefully(s.Engine)
}

func RunGracefully(mux http.Handler, fns ...func(context.Context) error) error {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}
	go func() {
		slog.Info(fmt.Sprintf("serving HTTP on :%s", port))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server ended abruptly: %s", "error", err.Error())
		} else {
			slog.Info("server ended gracefully")
		}

		stop()
	}()

	var wg sync.WaitGroup
	for _, fn := range fns {
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			err := fn(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "process ended abruptly", "error", err.Error())
			} else {
				slog.InfoContext(ctx, "process ended gracefully")
			}

			stop()
		}(ctx)
	}

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	wg.Wait()

	slog.Info("server exited")
	return nil
}
