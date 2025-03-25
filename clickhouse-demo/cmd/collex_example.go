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

// This example demonstrates how you would use collex with the ClickHouse exporter
// in a real application. Note that some parts are simulated since we're focusing on
// understanding the concepts without implementing the full adapter.

func ExampleWithCollex() {
	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Step 1: Create the ClickHouse exporter factory
	chFactory := clickhouseexporter.NewFactory()

	// Step 2: Wrap the factory with collex
	// In a real implementation, you would use:
	// collexFactory, err := collex.NewFactory(chFactory, nil)
	// if err != nil {
	//     log.Fatalf("Failed to create collex factory: %v", err)
	// }
	
	// Step 3: Configure the exporter
	config := chFactory.CreateDefaultConfig().(*clickhouseexporter.Config)
	config.Endpoint = "tcp://localhost:9000"
	config.Username = "default"
	config.Password = "password"
	config.Database = "otel"
	config.TracesTableName = "otel_traces"
	config.CreateSchema = true
	config.TTL = 72 * time.Hour

	// Step 4: Create a span exporter using collex
	// In a real implementation, you would use:
	// exporter, err := collexFactory.SpanExporter(ctx, config)
	// if err != nil {
	//     log.Fatalf("Failed to create span exporter: %v", err)
	// }
	//
	// For this example, we'll just use a simulated exporter
	exporter := &simulatedExporter{config: config}

	// Step 5: Create a resource for the tracer provider
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("clickhouse-demo-service"),
			semconv.ServiceVersion("0.1.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	// Step 6: Create a tracer provider with the ClickHouse exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown tracer provider: %v", err)
		}
	}()

	// Step 7: Set the global tracer provider
	otel.SetTracerProvider(tp)

	// Step 8: Get a tracer and create spans
	tracer := tp.Tracer("clickhouse-demo")
	for i := 0; i < 5; i++ {
		parentCtx, parentSpan := tracer.Start(
			ctx,
			"parent-operation",
			trace.WithAttributes(attribute.String("custom.attribute", "custom-value")),
		)

		// Create child spans
		for j := 0; j < 2; j++ {
			_, childSpan := tracer.Start(
				parentCtx,
				fmt.Sprintf("child-operation-%d", j),
				trace.WithAttributes(
					attribute.Int("child.number", j),
					attribute.Float64("random.value", rand.Float64()),
				),
			)
			// Simulate some work
			time.Sleep(10 * time.Millisecond)
			childSpan.End()
		}

		parentSpan.End()
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Traces have been sent to ClickHouse")
} 