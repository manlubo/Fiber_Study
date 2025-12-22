package observability

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

// TraceMiddleware는 Fiber 요청마다 Root Span을 생성한다.
func TraceMiddleware() fiber.Handler {

	tracer := otel.Tracer("http")

	return func(c *fiber.Ctx) error {

		ctx, span := tracer.Start(
			c.Context(),
			c.Method()+" "+c.Path(),
		)
		defer span.End()

		// Fiber Context에 OTel Context 주입
		c.SetUserContext(ctx)

		return c.Next()
	}
}
