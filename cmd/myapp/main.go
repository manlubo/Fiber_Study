package main

import (
	"study/internal/config"
	"study/internal/database"
	"study/internal/feature/auth"
	"study/internal/middleware"
	"study/internal/router"
	"study/pkg/dbmetrics"
	"study/pkg/log"
	"study/pkg/response"
	"study/pkg/util"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// .env 파일 로드
	if err := godotenv.Load(util.GetPath(".env")); err != nil {
		log.Error(".env 파일을 찾을 수 없습니다.", log.MapErr("error", err))
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error("설정 파일을 불러오는데 실패했습니다", log.MapErr("error", err))
	}

	app := fiber.New()
	app.Use(middleware.Cors(&cfg.Cors))
	app.Use(middleware.ApiMetrics())

	postgresdb, err := database.NewPostgres(&cfg.Postgres)
	if err != nil {
		log.Error("PostgresDB 연결에 실패했습니다", log.MapErr("error", err))
		return
	}
	log.Info("PostgresDB 연결 성공")
	metricsDB := dbmetrics.New(postgresdb)

	// JWT 환경설정 등록
	jwtService := auth.NewJwtService(&cfg.JWT)

	// 쿠키 환경설정 등록
	cookieService := auth.NewCookieService(&cfg.Cookie)
	authMiddleware := middleware.NewAuthMiddlewareConfig(cfg.Cookie.Name)

	// 라우터 등록
	router.Register(app, metricsDB, jwtService, cookieService, authMiddleware)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(response.OK("Application is running PORT : "+cfg.App.Port, nil))
	})
	log.Info("Application is running", log.MapStr("port", cfg.App.Port))

	app.Listen(":" + cfg.App.Port)
}
