package middleware

import (
	"time"

	"study/pkg/dbmetrics"
	"study/pkg/log"

	"github.com/gofiber/fiber/v2"
)

func ApiMetrics() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// /metrics는 제외
		if c.Path() == "/metrics" {
			return c.Next()
		}

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

		log.Info(
			"[API 성능 로그]",
			log.MapStr("api", api),
			log.MapInt64("time", elapsed),
			log.MapInt("queryCount", metrics.QueryCount),
			log.MapInt64("maxQueryTime", metrics.MaxTimeMs),
		)

		if metrics.MaxQuery != "" {
			log.Debug(metrics.MaxQuery)
		}

		return err
	}
}
