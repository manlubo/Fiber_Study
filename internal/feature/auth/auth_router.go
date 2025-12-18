package auth

import (
	"github.com/gofiber/fiber/v2"
)

type AuthRouter struct {
	handler *AuthHandler
}

func NewAuthRouter(handler *AuthHandler) *AuthRouter {
	return &AuthRouter{handler: handler}
}

func (r *AuthRouter) RegisterRoutes(
	open fiber.Router,
) {
	api := open.Group("/auth")

	api.Post("/signup", r.handler.SignUp)
	api.Post("/login", r.handler.Login)
	api.Post("/refresh", r.handler.Refresh)
	api.Post("/logout", r.handler.Logout)

}

func (r *AuthRouter) RegisterAuthRoutes(
	auth fiber.Router,
) {
	apiAuth := auth.Group("/auth")
	if apiAuth == nil {
		return
	}
}
