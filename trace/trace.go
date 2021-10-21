package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

func NewTracer(svcName string, exporter tracesdk.SpanExporter) trace.Tracer {
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
		)),
	)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otel.SetTracerProvider(tp)
	return tp.Tracer(``)
}

var Tracer trace.Tracer
var tracerOnce sync.Once

func InitTracerWithJaegerExporter(svcName string) error {
	// Create the Jaeger exporter
	// endpoint will be set by env OTEL_EXPORTER_JAEGER_ENDPOINT or jaeger.WithEndpoint(url)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint())
	if err != nil {
		return err
	}
	tracerOnce.Do(func() {
		Tracer = NewTracer(svcName, exp)
	})
	return nil
}
