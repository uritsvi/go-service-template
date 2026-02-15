# OpenTelemetry Module

Provides OpenTelemetry metrics and logging instrumentation for Go services.

## Setup

1. **Create a service-specific metrics struct**:

```go
package metrics

import (
	"internal/config"
	"internal/logger"
	"internal/otel"
	"go.opentelemetry.io/otel/metric"
)

var M *OtelMetrics

type OtelMetrics struct {
	meter metric.Meter
	RequestTotal       metric.Int64Counter
	RequestProcessed   metric.Int64Counter
	RequestProcessedMs metric.Int64Histogram
	InflightRequests   metric.Int64UpDownCounter
}

func init() {
	cfg := config.Get()
	M = &OtelMetrics{}
	if cfg.OtelEnabled {
		if err := M.Setup(cfg.OtelCollectorEndpoint, cfg.ServiceName); err != nil {
			logger.L.Fatalf("Failed to initialize metrics: %s", err)
		}
	} else {
		M.setupMetrics(otel.Meter())
	}
}

func (m *OtelMetrics) Setup(endpoint, serviceName string) error {
	meter, err := otel.SetupOtelMetrics(endpoint, serviceName)
	if err != nil {
		return err
	}
	m.meter = meter
	return m.setupMetrics(meter)
}

func (m *OtelMetrics) setupMetrics(meter metric.Meter) error {
	var err error
	m.RequestTotal, err = meter.Int64Counter("request_total", metric.WithDescription("Number of requests"), metric.WithUnit("count"))
	return err
}
```

2. **Logging** is automatically set up when using `internal/logger` if `OtelEnabled` is true.

## Usage

Access metrics through the global `M` variable:

```go
import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Increment a counter
metrics.M.RequestTotal.Add(ctx, 1, metric.WithAttributes(
	attribute.String("route", "/api/users"),
))

// Record a histogram value
metrics.M.RequestProcessedMs.Record(ctx, 150, metric.WithAttributes(
	attribute.String("status", "200"),
	attribute.String("outcome", "success"),
))

// Increment/decrement an up/down counter
metrics.M.InflightRequests.Add(ctx, 1)  // Increment
metrics.M.InflightRequests.Add(ctx, -1) // Decrement
```

### Adding Custom Metrics

Declare metrics in your struct and create them in `setupMetrics`:

```go
type OtelMetrics struct {
	meter metric.Meter
	CustomCounter      metric.Int64Counter
	CustomHistogram    metric.Int64Histogram
	CustomUpDownCounter metric.Int64UpDownCounter
}

func (m *OtelMetrics) setupMetrics(meter metric.Meter) error {
	var err error
	m.CustomCounter, err = meter.Int64Counter("custom_counter", metric.WithDescription("Description"), metric.WithUnit("count"))
	if err != nil {
		return err
	}
	m.CustomHistogram, err = meter.Int64Histogram("custom_duration_ms", metric.WithDescription("Duration"), metric.WithUnit("ms"))
	if err != nil {
		return err
	}
	m.CustomUpDownCounter, err = meter.Int64UpDownCounter("custom_inflight", metric.WithDescription("In-flight operations"), metric.WithUnit("count"))
	return err
}
```

Use custom metrics:

```go
metrics.M.CustomCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("status", "success")))
metrics.M.CustomHistogram.Record(ctx, 150, metric.WithAttributes(attribute.String("operation", "process_data")))
metrics.M.CustomUpDownCounter.Add(ctx, 1)  // Increment
metrics.M.CustomUpDownCounter.Add(ctx, -1) // Decrement
```

## Configuration

Metrics and logs are configured via your service's config:
- `OtelEnabled` - Enable/disable OpenTelemetry
- `OtelCollectorEndpoint` - OTLP collector endpoint (e.g., `localhost:4317`)
- `ServiceName` - Service name for resource attributes

When enabled, metrics and logs are automatically sent to the configured OTLP collector endpoint. Metrics are exported every 5 seconds.

## Shutdown

For graceful shutdown:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
otel.ShutdownOtelMetrics(ctx)
otel.ShutdownOtelLogger(ctx)
```
