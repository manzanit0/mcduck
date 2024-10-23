package xtrace

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Tracer() trace.Tracer {
	return otel.GetTracerProvider().Tracer("github.com/manzanit0/mcduck/pkg/xtrace")
}

func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer().Start(ctx, spanName, opts...)
}

func HydrateContext(ctx context.Context, traceID string, spanID string) context.Context {
	if traceID == "" || spanID == "" {
		return ctx
	}

	traceIDTyped := trace.TraceID([]byte(traceID))
	spanIDTyped := trace.SpanID([]byte(spanID))

	return trace.ContextWithRemoteSpanContext(
		context.Background(),
		trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceIDTyped,
			SpanID:     spanIDTyped,
			TraceFlags: trace.FlagsSampled,
		}),
	)
}

func GetSpan(ctx context.Context) (context.Context, trace.Span) {
	span := trace.SpanFromContext(ctx)
	return ctx, span
}

func RecordError(ctx context.Context, description string, err error) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, description)
	span.RecordError(err)

	slog.ErrorContext(ctx, description, "error", err.Error())
}
