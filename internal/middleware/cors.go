package middleware

import (
	"strings"
	"study/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(cfg *config.Cors) fiber.Handler {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			for _, o := range cfg.AllowOrigins {
				if o == origin {
					return true
				}
			}
			return false
		},
		AllowMethods:     strings.Join(cfg.AllowMethods, ","),
		AllowHeaders:     strings.Join(cfg.AllowHeaders, ","),
		AllowCredentials: cfg.AllowCredentials,
	})
}
