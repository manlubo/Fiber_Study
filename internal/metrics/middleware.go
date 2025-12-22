package metrics

import (
	"strconv"
	"study/pkg/dbmetrics"

	"github.com/gofiber/fiber/v2"
)

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// /metrics는 제외 (중요)
		if c.Path() == "/metrics" {
			return c.Next()
		}

		err := c.Next()

		status := c.Response().StatusCode()
		method := c.Method()

		// 라우트 패턴 기준 (/users/:id)
		path := "unknown"
		if c.Route() != nil {
			path = c.Route().Path
		}

		statusStr := strconv.Itoa(status)

		// 전체 요청 카운트
		dbmetrics.HttpRequestsTotal.WithLabelValues(
			method,
			path,
			statusStr,
		).Inc()

		// 에러만 따로 카운트
		if status >= 400 {
			dbmetrics.HttpErrorsTotal.WithLabelValues(
				method,
				path,
				statusStr,
			).Inc()
		}

		return err
	}
}
