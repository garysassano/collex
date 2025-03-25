# ClickHouse Exporter Demo with Collex

This demo shows how to use the OpenTelemetry Collector's ClickHouse exporter with the [collex](https://github.com/user/collex) library. The collex library allows you to use OpenTelemetry Collector exporters directly in Go applications, without requiring a separate OpenTelemetry Collector instance.

## Overview

The demo includes:

1. A simulated application that generates OpenTelemetry traces
2. Integration with the ClickHouse exporter using collex
3. A Docker Compose configuration to run ClickHouse locally

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose

## Running the Demo

### 1. Start ClickHouse

```bash
docker-compose up -d
```

This will start a ClickHouse server on port 9000 (native protocol) and 8123 (HTTP).

### 2. Run the Demo Application

```bash
# Run the simulated example (default)
go run cmd/main.go

# Run the more realistic implementation
go run cmd/main.go -example=real

# Show help
go run cmd/main.go -help
```

The demo offers two examples:

1. **Simulation Example (`-example=sim`)**: A simpler implementation that simulates what collex would do with the ClickHouse exporter. It generates spans continuously until interrupted with Ctrl+C.

2. **Real Implementation Example (`-example=real`)**: A more realistic implementation that's closer to how collex actually works with the ClickHouse exporter. It generates a fixed number of spans and then exits.

Both examples demonstrate the key steps for using collex with the ClickHouse exporter.

### Simulated Exporter Considerations

The demo includes a simulated version of what collex would actually do. In a real implementation:

1. collex would use the transmute package to convert OpenTelemetry SDK spans to the format expected by the ClickHouse exporter
2. The real exporter would establish a connection to ClickHouse and insert the spans into the database

For demonstration purposes, we've included detailed comments showing how you would use collex in a real application.

## Implementation Notes

Due to API compatibility issues, some parts of the code are commented out or simplified. In a real production application, you would:

1. Import the actual collex library
2. Properly handle connection errors and shutdown procedures
3. Use the appropriate version of dependencies that work together

## How It Works

The demo shows the key steps to use collex with the ClickHouse exporter:

1. Create a factory for the ClickHouse exporter
2. Wrap that factory with collex (simulated in this demo)
3. Configure the exporter with connection details for ClickHouse
4. Use the exporter with the OpenTelemetry Go SDK's `TracerProvider`
5. Generate spans that will be sent to ClickHouse

## Querying ClickHouse

After running the demo, you can query ClickHouse to see the traces that were inserted:

```sql
SELECT 
    TraceId,
    SpanId,
    ParentSpanId,
    SpanName,
    ServiceName,
    Duration,
    StatusCode,
    StatusMessage,
    toString(SpanAttributes) AS Attributes
FROM otel_traces
LIMIT 100;
```

You can run this query using the ClickHouse client or HTTP interface:

```bash
# Using curl with the HTTP interface
curl "http://localhost:8123/?query=SELECT+TraceId,SpanId,ParentSpanId,SpanName,ServiceName,Duration+FROM+otel_traces+LIMIT+10"

# Using the ClickHouse client (if installed)
clickhouse-client --host=localhost --port=9000 --user=default --password=password --query="SELECT TraceId, SpanId, ParentSpanId, SpanName, ServiceName, Duration FROM otel_traces LIMIT 10"
```

## Understanding the Code

- `main.go`: Demonstrates generating telemetry data with the OpenTelemetry SDK
- `collex_example.go`: Shows the steps to integrate collex with the ClickHouse exporter in a detailed, step-by-step manner

In a real implementation, you would use the actual collex library instead of the simulated version used in this demo. 