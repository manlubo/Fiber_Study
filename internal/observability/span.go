package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func StartServiceSpan(
	ctx context.Context,
	name string,
) (context.Context, trace.Span, time.Time) {
	start := time.Now()
	ctx, span := Tracer.Start(ctx, name)
	return ctx, span, start
}

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
