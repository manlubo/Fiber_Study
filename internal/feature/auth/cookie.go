package auth

import (
	"study/internal/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CookieService struct {
	Name     string
	Path     string
	HttpOnly bool
	Secure   bool
	SameSite string
	MaxAge   int
}

func NewCookieService(cfg *config.Cookie) *CookieService {
	return &CookieService{
		Name:     cfg.Name,
		Path:     cfg.Path,
		HttpOnly: cfg.HttpOnly,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		MaxAge:   cfg.MaxAge,
	}
}

// 쿠키 생성
func (s *CookieService) SetCookie(c *fiber.Ctx, refreshToken string, rememberMe bool) error {
	cookie := &fiber.Cookie{
		Name:     s.Name,
		Value:    refreshToken,
		Path:     s.Path,
		HTTPOnly: s.HttpOnly,
		Secure:   s.Secure,
		SameSite: s.SameSite,
	}

	if rememberMe {
		cookie.MaxAge = s.MaxAge * 60 * 60 * 24
	} else {
		cookie.MaxAge = 0
		cookie.Expires = time.Time{}
	}

	c.Cookie(cookie)

	return c.SendStatus(fiber.StatusOK)
}

// 쿠키 삭제
func (s *CookieService) RemoveCookie(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     s.Name,
		Value:    "",
		Path:     s.Path,
		HTTPOnly: s.HttpOnly,
		Secure:   s.Secure,
		SameSite: s.SameSite,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	return c.SendStatus(fiber.StatusOK)
}

// 쿠키 조회
func (s *CookieService) GetCookie(c *fiber.Ctx) string {
	return c.Cookies(s.Name)
}
