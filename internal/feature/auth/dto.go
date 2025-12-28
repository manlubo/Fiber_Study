package auth

import (
	"study/internal/feature/member"
)

// 회원가입 요청 DTO
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// 로그인 요청 DTO
type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberMe" default:"false"`
}

// 멤버 전달 객체
type MemberResponse struct {
	ID      int64         `json:"id"`
	Email   string        `json:"email"`
	Profile *string       `json:"profile"`
	Roles   []member.Role `json:"roles"`
}

// 로그인 응답 DTO
type LoginResponse struct {
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	Member       MemberResponse `json:"member"`
}
