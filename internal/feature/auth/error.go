package auth

import "errors"

// 서비스 에러
var (
	// 이미 존재하는 이메일
	ErrEmailAlreadyExists = errors.New("EMAIL_ALREADY_EXISTS")

	// 이메일 또는 비밀번호가 일치하지 않음
	ErrInvalidCredential = errors.New("INVALID_CREDENTIAL")

	// 토큰 만료
	ErrTokenExpired = errors.New("TOKEN_EXPIRED")

	// 토큰 유효하지 않음
	ErrTokenInvalid = errors.New("TOKEN_INVALID")

	// 토큰 타입 불일치
	ErrTokenTypeWrong = errors.New("TOKEN_TYPE_WRONG")

	// 쿠키 누락
	ErrCookieNotFound = errors.New("COOKIE_NOT_FOUND")
)
