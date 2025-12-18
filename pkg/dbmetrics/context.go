package dbmetrics

import "context"

type ctxKey struct{}

func NewContext(ctx context.Context) (context.Context, *Metrics) {
	m := &Metrics{}
	return context.WithValue(ctx, ctxKey{}, m), m
}

func FromContext(ctx context.Context) *Metrics {
	m, _ := ctx.Value(ctxKey{}).(*Metrics)
	return m
}
