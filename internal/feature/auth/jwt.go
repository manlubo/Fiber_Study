package auth

import (
	"errors"
	"study/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT(access / refresh) 생성 및 검증을 담당하는 서비스
type JwtService struct {
	accessSecret     []byte
	refreshSecret    []byte
	accessExpireMin  int
	refreshExpireDay int
}

// JWT 설정값을 기반으로 JwtService 생성
func NewJwtService(cfg *config.JWT) *JwtService {
	return &JwtService{
		accessSecret:     cfg.AccessSecret,
		refreshSecret:    cfg.RefreshSecret,
		accessExpireMin:  cfg.AccessExpireMin,
		refreshExpireDay: cfg.RefreshExpireDay,
	}
}

// JWT Payload에 담기는 공통 클레임 구조
type Claims struct {
	MemberID int64     `json:"memberId"`
	Type     TokenType `json:"type"`
	jwt.RegisteredClaims
}

// 토큰 종류 구분 (Access / Refresh)
type TokenType string

const (
	TypeAccess  TokenType = "ACCESS"
	TypeRefresh TokenType = "REFRESH"
)

// access 토큰 상태 분류 (미들웨어용)
type TokenStatus int

const (
	TokenValid TokenStatus = iota
	TokenExpired
	TokenInvalid
)

// Access Token 생성
func (j *JwtService) GenerateAccessToken(memberID int64) (string, error) {
	return j.generateToken(memberID, TypeAccess)
}

// Refresh Token 생성
func (j *JwtService) GenerateRefreshToken(memberID int64) (string, error) {
	return j.generateToken(memberID, TypeRefresh)
}

// 로그인
func (j *JwtService) Login(memberID int64) (*LoginResponse, error) {
	accessToken, err := j.GenerateAccessToken(memberID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := j.GenerateRefreshToken(memberID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// access 토큰 검증 및 Claims 반환
func (j *JwtService) VerifyAccessToken(tokenStr string) (*Claims, error) {
	return j.verifyToken(tokenStr, TypeAccess)
}

// refresh 토큰 검증 및 Claims 반환
func (j *JwtService) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return j.verifyToken(tokenStr, TypeRefresh)
}

// 토큰 생성 공통 로직 (access / refresh)
func (j *JwtService) generateToken(memberID int64, tokenType TokenType) (string, error) {
	var (
		secret     []byte
		expireTime time.Time
	)

	now := time.Now()

	switch tokenType {
	case TypeAccess:
		secret = j.accessSecret
		expireTime = now.Add(time.Duration(j.accessExpireMin) * time.Minute)

	case TypeRefresh:
		secret = j.refreshSecret
		expireTime = now.Add(time.Duration(j.refreshExpireDay) * 24 * time.Hour)

	default:
		return "", jwt.ErrTokenInvalidClaims
	}

	claims := Claims{
		MemberID: memberID,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(secret)
}

// 토큰 검증 공통 로직 (서명, 만료, 타입 검증)
func (j *JwtService) verifyToken(tokenStr string, tokenType TokenType) (*Claims, error) {
	var secret []byte

	switch tokenType {
	case TypeAccess:
		secret = j.accessSecret
	case TypeRefresh:
		secret = j.refreshSecret
	default:
		return nil, ErrTokenInvalid
	}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			// 서명 알고리즘 검증
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			return secret, nil
		},
	)

	// 파싱 에러 처리 (만료 / 위조)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// Claims 타입 체크
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	// JWT 유효성
	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	// 토큰 타입 검증
	if claims.Type != tokenType {
		return nil, ErrTokenTypeWrong
	}

	return claims, nil
}

// Authorization 헤더에서 Bearer 토큰 추출
func ExtractBearer(authHeader string) string {
	const prefix = "Bearer "

	if authHeader == "" {
		return ""
	}
	if len(authHeader) <= len(prefix) {
		return ""
	}
	if authHeader[:len(prefix)] != prefix {
		return ""
	}
	return authHeader[len(prefix):]
}

// access 토큰 상태를 Valid / Expired / Invalid 로 판별
func (j *JwtService) VerifyStatus(tokenStr string) TokenStatus {
	_, err := j.VerifyAccessToken(tokenStr)

	if err == nil {
		return TokenValid
	}

	switch err {
	case ErrTokenExpired:
		return TokenExpired
	default:
		return TokenInvalid
	}
}

// 검증된 access 토큰에서 Claims만 추출
func (j *JwtService) Claims(tokenStr string) *Claims {
	claims, err := j.VerifyAccessToken(tokenStr)
	if err != nil {
		return nil
	}
	return claims
}

// refresh 토큰 유효성만 검증 (미들웨어용)
func (j *JwtService) Verify(tokenStr string) error {
	_, err := j.VerifyRefreshToken(tokenStr)
	return err
}
