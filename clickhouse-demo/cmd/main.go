package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create ClickHouse config
	cfg := &clickhouseexporter.Config{
		Endpoint:        "tcp://localhost:9000",
		Username:        "default",
		Password:        "password",
		Database:        "otel",
		TracesTableName: "otel_traces",
		CreateSchema:    true,
		TTL:             72 * time.Hour,
	}

	// Create collex factory for ClickHouse
	// In a real implementation, you would use:
	// factory, err := collex.NewFactory(clickhouseexporter.NewFactory(), nil)
	// 
	// But for this example, we'll create a simplified version since we're focusing on understanding
	// how the exporter works with the collector

	// Create resource with identifying information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("clickhouse-demo-service"),
			semconv.ServiceVersion("0.1.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// Create a trace exporter
	// In a real implementation with collex, you would use:
	// exp, err := factory.SpanExporter(ctx, cfg)
	// Here we're just simulating what this would do
	traceExporter := &simulatedExporter{
		config: cfg,
	}

	// Create tracer provider with the ClickHouse exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown tracer provider: %v", err)
		}
	}()

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Get a tracer
	tracer := tp.Tracer("clickhouse-demo")

	fmt.Println("Starting to generate telemetry data. Press Ctrl+C to stop.")
	fmt.Println("This demo shows how you would use collex with the ClickHouse exporter.")
	fmt.Println("In a real implementation, collex provides an adapter between")
	fmt.Println("the OpenTelemetry Collector exporters and the OpenTelemetry Go SDK.")
	fmt.Println("\nClickHouse Configuration:")
	fmt.Printf("  Endpoint: %s\n", cfg.Endpoint)
	fmt.Printf("  Database: %s\n", cfg.Database) 
	fmt.Printf("  Traces Table: %s\n", cfg.TracesTableName)
	fmt.Println("\nPress Ctrl+C to stop.")

	// Generate traces every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down...")
			return
		case <-ticker.C:
			// Create a parent span
			parentCtx, parentSpan := tracer.Start(
				ctx,
				"parent-operation",
				trace.WithAttributes(attribute.String("custom.attribute", "custom-value")),
			)

			// Generate a random number of child spans (1-3)
			numChildSpans := rand.Intn(3) + 1
			for i := 0; i < numChildSpans; i++ {
				_, childSpan := tracer.Start(
					parentCtx,
					fmt.Sprintf("child-operation-%d", i),
					trace.WithAttributes(
						attribute.Int("child.number", i),
						attribute.Float64("random.value", rand.Float64()),
					),
				)
				// Simulate some work being done
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
				childSpan.End()
			}

			// End the parent span
			parentSpan.End()
			fmt.Printf("Generated trace with %d child spans\n", numChildSpans)
		}
	}
}

// simulatedExporter simulates what collex would do with the ClickHouse exporter
type simulatedExporter struct {
	config *clickhouseexporter.Config
}

func (e *simulatedExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	// In a real implementation with collex, this would convert the spans to collector format
	// and send them to ClickHouse using the exporter's functionality
	fmt.Printf("Exporting %d spans to ClickHouse at %s\n", len(spans), e.config.Endpoint)
	return nil
}

func (e *simulatedExporter) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down ClickHouse exporter")
	return nil
} 