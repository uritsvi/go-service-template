package otel

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

var (
	otelLoggerProvider *sdklog.LoggerProvider
	otelMeterProvider  *sdkmetric.MeterProvider
	otelMeter          metric.Meter
	metricsLock        sync.RWMutex
)

type OtelHook struct {
	logger otellog.Logger
}

func (h *OtelHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *OtelHook) Fire(entry *logrus.Entry) error {
	var severity otellog.Severity
	switch entry.Level {
	case logrus.PanicLevel, logrus.FatalLevel:
		severity = otellog.SeverityFatal
	case logrus.ErrorLevel:
		severity = otellog.SeverityError
	case logrus.WarnLevel:
		severity = otellog.SeverityWarn
	case logrus.InfoLevel:
		severity = otellog.SeverityInfo
	case logrus.DebugLevel, logrus.TraceLevel:
		severity = otellog.SeverityDebug
	default:
		severity = otellog.SeverityInfo
	}

	attrs := make([]otellog.KeyValue, 0, len(entry.Data)+3)

	for k, v := range entry.Data {
		attrs = append(attrs, otellog.String(k, fmt.Sprintf("%v", v)))
	}

	record := otellog.Record{}
	record.SetTimestamp(entry.Time)
	record.SetSeverity(severity)
	record.SetBody(otellog.StringValue(entry.Message))
	record.AddAttributes(attrs...)

	h.logger.Emit(context.Background(), record)

	return nil
}

func SetupOtelLogger(endpoint, serviceName string, logger *logrus.Logger) error {
	ctx := context.Background()

	exporter, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithReconnectionPeriod(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	batchProcessor := sdklog.NewBatchProcessor(exporter,
		sdklog.WithExportMaxBatchSize(512),
		sdklog.WithExportInterval(5*time.Second),
		sdklog.WithExportTimeout(30*time.Second),
		sdklog.WithMaxQueueSize(2048),
		sdklog.WithExportBufferSize(512),
	)

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(batchProcessor),
	)

	otelLogger := loggerProvider.Logger("go/logger")

	hook := &OtelHook{logger: otelLogger}
	logger.AddHook(hook)

	otelLoggerProvider = loggerProvider

	return nil
}

func ShutdownOtelLogger(ctx context.Context) error {
	if otelLoggerProvider != nil {
		return otelLoggerProvider.Shutdown(ctx)
	}

	return nil
}

func SetupOtelMetrics(endpoint, serviceName string) (metric.Meter, error) {
	metricsLock.Lock()
	defer metricsLock.Unlock()

	ctx := context.Background()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithReconnectionPeriod(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter,
		sdkmetric.WithInterval(5*time.Second),
	)

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	otel.SetMeterProvider(provider)

	meter := provider.Meter(serviceName)

	otelMeterProvider = provider
	otelMeter = meter

	return meter, nil
}

func ShutdownOtelMetrics(ctx context.Context) error {
	metricsLock.Lock()
	defer metricsLock.Unlock()

	if otelMeterProvider != nil {
		return otelMeterProvider.Shutdown(ctx)
	}

	return nil
}

func Meter() metric.Meter {
	if otelMeter == nil {
		otelMeter = otel.Meter("unknown")
	}

	return otelMeter
}
