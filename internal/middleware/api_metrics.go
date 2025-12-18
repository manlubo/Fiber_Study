// internal/middleware/api_metrics.go
package middleware

import (
	"time"

	"study/pkg/dbmetrics"
	"study/pkg/log"

	"github.com/gofiber/fiber/v2"
)

func ApiMetrics() fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, metrics := dbmetrics.NewContext(c.UserContext())
		c.SetUserContext(ctx)

		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start).Milliseconds()

		log.Info(
			"[API 성능 로그]",
			log.MapStr("api", c.Route().Path),
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
