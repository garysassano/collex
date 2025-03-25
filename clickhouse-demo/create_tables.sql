CREATE TABLE IF NOT EXISTS otel.otel_traces (
    Timestamp DateTime64(9) DEFAULT now(),
    TraceId String,
    SpanId String,
    ParentSpanId String,
    SpanName String,
    SpanKind String,
    ServiceName String,
    ResourceAttributes String,
    SpanAttributes String,
    Duration UInt64,
    StatusCode UInt16,
    StatusMessage String,
    Events String,
    Links String
) ENGINE = MergeTree()
ORDER BY (ServiceName, SpanName, toDate(Timestamp))
PARTITION BY toYYYYMM(Timestamp); 