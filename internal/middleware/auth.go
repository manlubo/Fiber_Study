package middleware

import (
	"study/internal/feature/auth"
	"study/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddlewareConfig struct {
	CookieName string
}

func NewAuthMiddlewareConfig(cookieName string) *AuthMiddlewareConfig {
	return &AuthMiddlewareConfig{CookieName: cookieName}
}

func (cfg *AuthMiddlewareConfig) AuthMiddleware(jwtSvc *auth.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Authorization 헤더에서 Bearer 토큰 추출
		access := auth.ExtractBearer(c.Get("Authorization"))

		// refresh 토큰은 HttpOnly 쿠키에서만 읽음
		refresh := c.Cookies(cfg.CookieName)

		// access token이 있는 경우
		if access != "" {

			// access token 상태 판별
			// Valid   : 정상
			// Expired : 만료
			// Invalid : 위조 / 서명 불일치 / 형식 오류
			status := jwtSvc.VerifyStatus(access)

			switch status {

			case auth.TokenValid:
				// 정상 토큰인 경우
				claims := jwtSvc.Claims(access)
				if claims == nil {
					return c.Status(401).JSON(response.Error("INVALID_TOKEN", "Invalid token", nil))
				}
				c.Locals("claims", claims)
				return c.Next()

			case auth.TokenExpired:
				// access 토큰은 만료되었지만, refresh 토큰이 있고 유효
				if refresh != "" && jwtSvc.Verify(refresh) == nil {
					return c.Status(401).JSON(response.Error("ACCESS_EXPIRED", "Access token expired", nil))
				}

				// refresh 토큰도 없거나 만료
				return c.Status(401).JSON(response.Error("SESSION_EXPIRED", "Session expired", nil))

			case auth.TokenInvalid:
				// access 토큰이 위조되었거나 조작됨
				return c.Status(401).JSON(response.Error("INVALID_TOKEN", "Invalid token", nil))
			}
		}

		// access는 없지만 refresh가 있고 유효한 경우
		if refresh != "" && jwtSvc.Verify(refresh) == nil {
			return c.Status(401).JSON(response.Error("ACCESS_REQUIRED", "Access token required", nil))
		}

		// access도 없고 refresh도 없거나 만료
		return c.Status(401).JSON(response.Error("SESSION_EXPIRED", "Session expired", nil))
	}
}
