package router

import (
	"study/internal/feature/auth"
	"study/internal/feature/member"
	"study/internal/middleware"
	"study/internal/shared/db"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, db db.DB, jwtService *auth.JwtService, cookieService *auth.CookieService, authMiddleware *middleware.AuthMiddlewareConfig) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// member
	memberRepo := member.NewMemberRepository(db)

	// auth
	authService := auth.NewAuthService(memberRepo, jwtService)
	authHandler := auth.NewAuthHandler(authService, cookieService)
	authRouter := auth.NewAuthRouter(authHandler)

	// ==================================== 인증 필요 없음
	authRouter.RegisterRoutes(v1)

	// ==================================== 인증 필요
	v1Auth := v1.Group("", authMiddleware.AuthMiddleware(jwtService))
	authRouter.RegisterAuthRoutes(v1Auth)

}
