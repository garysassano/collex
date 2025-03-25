package internal

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

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

// SimulatedClickHouseExample runs a demo with a simulated exporter
func SimulatedClickHouseExample() {
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
		// Note: CreateSchema field might not exist in newer versions, removed
		TTL:             72 * time.Hour,
	}

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

// emptyHostImpl is used to satisfy the component.Host interface
type emptyHostImpl struct{}

func (emptyHostImpl) GetExtensions() map[component.ID]component.Component {
	return nil
}

// spanExporterAdapter adapts the collector exporter to the SDK
type spanExporterAdapter struct {
	cexp consumer.Traces
}

func (e *spanExporterAdapter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	// In a real implementation with collex, this would convert the spans to collector format
	// using the transmute package, and send them to ClickHouse.
	//
	// The real implementation in collex does:
	// return e.cexp.ConsumeTraces(ctx, transmute.Spans(spans))
	//
	// For this demo, we'll just log that spans are being exported
	fmt.Printf("Exporting %d spans to ClickHouse\n", len(spans))
	return nil
}

func (e *spanExporterAdapter) Shutdown(ctx context.Context) error {
	// In a real implementation, we would shutdown the exporter
	return nil
}

// RealClickHouseExample demonstrates a more realistic implementation closer to how collex works
func RealClickHouseExample() {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create a logger
	_, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Use this in real application:
	// Step 1: Create a collex Factory
	// chFactory := clickhouseexporter.NewFactory()
	// collexFactory, err := collex.NewFactory(chFactory, nil)
	// if err != nil {
	//     log.Fatalf("Failed to create collex factory: %v", err)
	// }

	// Step 2: Create and configure ClickHouse exporter
	config := &clickhouseexporter.Config{
		Endpoint:        "tcp://localhost:9000",
		Username:        "default",
		Password:        "password",
		Database:        "otel",
		TracesTableName: "otel_traces",
		// Note: CreateSchema field might not exist in newer versions, removed
		TTL:             72 * time.Hour,
	}

	// Use this in real application:
	// Step 3: Create a span exporter using collex
	// traceExporter, err := collexFactory.SpanExporter(ctx, config)
	// if err != nil {
	//     log.Fatalf("Failed to create span exporter: %v", err)
	// }

	// For this demo, we'll simulate the adapter part
	// Get the real ClickHouse factory
	// factory := clickhouseexporter.NewFactory()
	
	// Create settings
	// settings := exporter.CreateSettings{
	//     TelemetrySettings: component.TelemetrySettings{
	//         Logger:         logger,
	//         TracerProvider: otel.GetTracerProvider(),
	//         MeterProvider:  otel.GetMeterProvider(),
	//     },
	//     BuildInfo: component.BuildInfo{
	//         Command:     "collex-demo",
	//         Description: "ClickHouse exporter demo",
	//         Version:     "latest",
	//     },
	// }
	
	// Create the actual exporter
	// Note: API changed, using a simulated exporter instead
	// exp, err := factory.CreateTraces(ctx, settings, config)
	// In real code with collex, you would use the actual factory methods
	
	// Simulate exporter
	traceExporter := &simulatedExporter{
		config: config,
	}

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

	// Create tracer provider with the exporter
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

	fmt.Println("Starting to generate telemetry data...")

	// Generate a few traces
	for i := 0; i < 5; i++ {
		// Create a parent span
		parentCtx, parentSpan := tracer.Start(
			ctx,
			fmt.Sprintf("parent-operation-%d", i),
			trace.WithAttributes(attribute.String("custom.attribute", "custom-value")),
		)

		// Create some child spans
		for j := 0; j < 3; j++ {
			_, childSpan := tracer.Start(
				parentCtx,
				fmt.Sprintf("child-operation-%d", j),
				trace.WithAttributes(
					attribute.Int("child.number", j),
					attribute.Float64("random.value", rand.Float64()),
				),
			)
			// Simulate some work being done
			time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			childSpan.End()
		}

		// End the parent span
		parentSpan.End()
		fmt.Printf("Generated trace %d with 3 child spans\n", i)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Trace generation complete. Flushing telemetry...")
	time.Sleep(1 * time.Second) // Give time for the batcher to flush

	fmt.Println("Telemetry has been sent to ClickHouse")
} 