package test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTraceSpanCreation(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	tracer := otel.Tracer("test-tracer")
	_, span := tracer.Start(context.Background(), "TestTraceSpan")
	span.End()

	// Ensure exporter captured the span
	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("no spans created")
	}
	if spans[0].Name != "TestTraceSpan" {
		t.Errorf("expected span name 'TestTraceSpan', got %s", spans[0].Name)
	}
}
