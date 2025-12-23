package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// span 시작
func StartServiceSpan(
	ctx context.Context,
	name string,
) (context.Context, trace.Span, time.Time) {
	start := time.Now()
	ctx, span := Tracer.Start(ctx, "service."+name)
	return ctx, span, start
}

// span 종료
func EndSpanWithLatency(
	span trace.Span,
	start time.Time,
	thresholdMs int64,
) {
	elapsed := time.Since(start)

	if elapsed.Milliseconds() > thresholdMs {
		span.SetAttributes(
			attribute.Bool("slow", true),
			attribute.Int64("elapsed_ms", elapsed.Milliseconds()),
			attribute.Int64("slow_threshold_ms", thresholdMs),
		)
	}

	span.End()
}

// 에러 기록
func RecordServiceError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// 커스텀 에러 기록(비즈니스 에러)
func RecordBusinessError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.AddEvent(
		"business.error",
		trace.WithAttributes(
			attribute.String("error.code", err.Error()),
		),
	)
}
