package observability

import (
	"context"
	"study/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

// InitTracer는 OpenTelemetry TracerProvider를 초기화하고
// 종료 시 호출할 shutdown 함수를 반환한다.
func InitTracer(cfg *config.Observability) (func(context.Context) error, error) {

	// disable 스위치
	if !cfg.Enabled {
		otel.SetTracerProvider(sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.NeverSample()),
		))
		return func(context.Context) error { return nil }, nil
	}

	// OTLP Exporter
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(cfg.OtlpEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Sampling 적용 (ParentBased + Ratio)
	sampler := sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(cfg.SampleRatio),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.ServiceName),
			),
		),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
