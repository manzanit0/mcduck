package xtrace

import (
	"github.com/gin-gonic/gin"
	"github.com/manzanit0/mcduck/pkg/auth"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func GinTraceRequests(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}

func GinEnhanceTraceAttributes() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		span.SetAttributes(attribute.String("mduck.user.email", auth.GetUserEmail(c)))
		c.Next()
	}
}
