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
go run cmd/main.go
```

This will generate spans and send them to ClickHouse. The application will continue running and generating spans until interrupted (Ctrl+C).

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