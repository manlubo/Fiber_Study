package auth

import (
	"study/internal/shared/errorx"
	"study/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// Handler
type AuthHandler struct {
	service       *AuthService
	cookieService *CookieService
}

func NewAuthHandler(service *AuthService, cookieService *CookieService) *AuthHandler {
	return &AuthHandler{service: service, cookieService: cookieService}
}

// 회원가입
func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req SignUpRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errorx.ErrRequestParseFailed.Error(), "JSON 파싱 실패", nil))
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errorx.ErrRequiredFieldMissing.Error(), "필수값 누락", nil))
	}

	if err := h.service.Register(ctx, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error(), "회원가입 실패", nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.OK("회원가입 성공", nil))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errorx.ErrRequestParseFailed.Error(), "JSON 파싱 실패", nil))
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errorx.ErrRequiredFieldMissing.Error(), "필수값 누락", nil))
	}

	loginResponse, err := h.service.Login(ctx, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error(), "로그인 실패", nil))
	}
	// 쿠키 생성
	_ = h.cookieService.SetCookie(c, loginResponse.RefreshToken, req.RememberMe)
	loginResponse.RefreshToken = ""

	return c.Status(fiber.StatusOK).JSON(response.OK("로그인 성공", loginResponse))
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	ctx := c.UserContext()

	refreshToken := h.cookieService.GetCookie(c)
	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(ErrCookieNotFound.Error(), "리프레쉬 쿠키 누락", nil))
	}
	loginResponse, err := h.service.Refresh(ctx, refreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error(), "리프레쉬 토큰으로 로그인 실패", nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.OK("리프레쉬 토큰으로 로그인 성공", loginResponse))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	h.cookieService.RemoveCookie(c)
	return c.Status(fiber.StatusOK).JSON(response.OK("로그아웃 성공", nil))
}
