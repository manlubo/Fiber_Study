package auth

import (
	"study/internal/config"
	"study/pkg/log"
	"testing"
)

/*
JWT_ACCESS_SECRET=ACCESS_SECRET_KEY
JWT_REFRESH_SECRET=REFRESH_SECRET_KEY
JWT_ACCESS_EXPIRE_MIN=30
JWT_REFRESH_EXPIRE_DAY=14
*/
func TestGenerateAccessToken(t *testing.T) {
	t.Setenv("APP_ENV", "dev")

	cfg, err := config.Load()
	if err != nil {
		log.Error("설정 파일을 불러오는데 실패했습니다", log.MapErr("error", err))
	}

	jwtService := NewJwtService(&cfg.JWT)
	memberId := int64(1)
	accessToken, err := jwtService.GenerateAccessToken(memberId)
	if err != nil {
		log.Error("Access Token 생성 실패", log.MapErr("error", err))
	}
	log.Info("Access Token 생성 성공", log.MapStr("accessToken", accessToken))
}
