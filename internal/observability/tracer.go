package observability

import "go.opentelemetry.io/otel"

var Tracer = otel.Tracer("myapp")
