package router

import (
	"study/internal/feature/auth"
	"study/internal/middleware"
	"study/internal/query"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(app *fiber.App, pool *pgxpool.Pool, queries *query.Queries, jwtService *auth.JwtService, cookieService *auth.CookieService, authMiddleware *middleware.AuthMiddlewareConfig) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// auth
	authService := auth.NewAuthService(pool, queries, jwtService)
	authHandler := auth.NewAuthHandler(authService, cookieService)
	authRouter := auth.NewAuthRouter(authHandler)

	// ==================================== 인증 필요 없음
	authRouter.RegisterRoutes(v1)

	// ==================================== 인증 필요
	v1Auth := v1.Group("", authMiddleware.AuthMiddleware(jwtService))
	authRouter.RegisterAuthRoutes(v1Auth)

}
