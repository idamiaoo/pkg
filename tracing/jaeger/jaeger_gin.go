package jaeger

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

type conetxtKey struct{}

var (
	traceIDKey = conetxtKey{}
)

// TraceID 从context里获取jaeger的TraceID
func TraceID(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			return sc.TraceID().String()
		}
	}
	return ""
}

// ContextWithTraceID 返回一个包含 traceID的 context
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// GinMiddleware gin中间件
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		spanOption := []opentracing.StartSpanOption{
			ext.SpanKindRPCServer,
		}
		spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if spanCtx != nil {
			spanOption = append(spanOption, opentracing.ChildOf(spanCtx))
		}
		span, ctx := opentracing.StartSpanFromContext(
			c.Request.Context(),
			c.Request.URL.Path,
			spanOption...,
		)
		defer span.Finish()

		ext.Component.Set(span, "HTTP")
		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		if sc, ok := span.Context().(jaeger.SpanContext); ok {
			c.Writer.Header().Set("X-Request-ID", sc.TraceID().String())
			ctx = ContextWithTraceID(ctx, sc.TraceID().String())
		}
		c.Request = c.Request.Clone(ctx)

		c.Next()
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		ext.Error.Set(span, c.Writer.Status() <= 199 || c.Writer.Status() >= 300)
	}
}
