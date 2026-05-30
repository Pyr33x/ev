package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/pyr33x/ev/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type Adapter struct {
	tracerProvider *tracesdk.TracerProvider
	tracer         trace.Tracer
}

func New(ctx context.Context, cfg *config.Telemetry) (*Adapter, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var exporter tracesdk.SpanExporter
	var err error

	if cfg.UseGRPC {
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		exporter, err = otlptracegrpc.New(ctx, opts...)
	} else {
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		exporter, err = otlptracehttp.New(ctx, opts...)
	}

	if err != nil {
		return nil, fmt.Errorf("create OTLP exporter: %w", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
		semconv.DeploymentEnvironment(cfg.Environment),
	)

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Adapter{
		tracerProvider: tp,
		tracer:         tp.Tracer(cfg.ServiceName),
	}, nil
}

func (a *Adapter) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, func()) {
	ctx, span := a.tracer.Start(ctx, name, opts...)
	return ctx, func() { span.End() }
}

func (a *Adapter) Tracer() trace.Tracer {
	return a.tracer
}

func (a *Adapter) Shutdown(ctx context.Context) error {
	if a.tracerProvider == nil {
		return nil
	}
	return a.tracerProvider.Shutdown(ctx)
}
