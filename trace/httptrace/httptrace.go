package httptrace

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	kittrace "../trace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/http/httptrace"
)

func NewGinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		r := c.Request
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		var span trace.Span
		ctx, span = kittrace.Tracer.Start(ctx,
			defaultOpNameFunc(r),
			trace.WithSpanKind(trace.SpanKindServer))
		c.Request = c.Request.WithContext(ctx)
		defer span.End()
		c.Next()
		span.SetAttributes(
			semconv.HTTPMethodKey.String(r.Method),
			semconv.HTTPSchemeKey.String(r.URL.Scheme),
			semconv.HTTPStatusCodeKey.Int(c.Writer.Status()),
		)
	}
}

// defaultOpNameFunc is default function that get operation name from http request
func defaultOpNameFunc(r *http.Request) string {
	return r.URL.Scheme + " " + r.Method + " " + r.URL.Path
}

type Client struct {
	*http.Client
}

func NewClient() Client {
	return Client{
		Client: &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)},
	}
}

// make sure request is new with http.NewRequestWithContext()
func (c Client) Do(r *http.Request) (*http.Response, error) {
	if !(r.Context() == context.Background() || r.Context() == context.TODO()) {
		ctx, span := kittrace.Tracer.Start(r.Context(), fmt.Sprintf(`%s: %s`, r.Method, r.URL))
		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
		r = r.WithContext(ctx)
		defer span.End()
	}

	return c.Client.Do(r)
}
