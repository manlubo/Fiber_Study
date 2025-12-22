package middleware

import (
	"time"

	"study/pkg/dbmetrics"

	"github.com/gofiber/fiber/v2"
)

func ApiMetrics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// /metrics는 제외
		if c.Path() == "/metrics" {
			return c.Next()
		}

		// DB metrics context 주입
		ctx, metrics := dbmetrics.NewContext(c.UserContext())
		c.SetUserContext(ctx)

		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start).Milliseconds()

		// api path 안전하게 가져오기
		api := "unknown"
		if c.Route() != nil {
			api = c.Route().Path
		}

		// Prometheus Histogram (DB 최대 쿼리 시간)
		if metrics.MaxTimeMs > 0 {
			dbmetrics.QueryMaxDuration.
				WithLabelValues(api).
				Observe(float64(metrics.MaxTimeMs) / 1000) // ms → sec
		}

		// API 응답 시간 통계
		dbmetrics.ApiDuration.
			WithLabelValues(api).
			Observe(float64(elapsed) / 1000)

		return err
	}
}
