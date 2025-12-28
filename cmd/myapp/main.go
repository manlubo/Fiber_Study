package main

import (
	"context"
	"study/internal/config"
	"study/internal/database"
	"study/internal/feature/auth"
	"study/internal/metrics"
	"study/internal/middleware"
	"study/internal/observability"
	"study/internal/query"
	"study/internal/router"
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

	// config 로드
	cfg, err := config.Load()
	if err != nil {
		log.Error("설정 파일을 불러오는데 실패했습니다", log.MapErr("error", err))
	}

	// OpenTelemetry 초기화
	shutdown, err := observability.InitTracer(&cfg.Observability)
	if err != nil {
		log.Error("OpenTelemetry 초기화에 실패했습니다", log.MapErr("error", err))
		return
	}
	defer shutdown(context.Background())

	// fiber app 생성
	app := fiber.New()

	// metrics 초기화
	metrics.Init()

	// middleware 등록
	app.Use(observability.TraceMiddleware())
	app.Use(middleware.Cors(&cfg.Cors))
	app.Use(metrics.Middleware())
	app.Use(middleware.ApiMetrics())

	// db 연결
	postgresdb, err := database.NewPostgres(&cfg.Postgres)
	if err != nil {
		log.Error("PostgresDB 연결에 실패했습니다", log.MapErr("error", err))
		return
	}
	log.Info("PostgresDB 연결 성공")
	if err := database.RunMigration(database.CreateDsn(&cfg.Postgres)); err != nil {
		log.Error("데이터베이스 마이그레이션에 실패했습니다", log.MapErr("error", err))
		return
	}
	log.Info("데이터베이스 마이그레이션 성공")

	queries := query.New(postgresdb)

	// auth 관련
	jwtService := auth.NewJwtService(&cfg.JWT)
	cookieService := auth.NewCookieService(&cfg.Cookie)
	authMiddleware := middleware.NewAuthMiddlewareConfig(cfg.Cookie.Name)

	// 라우터
	router.Register(app, postgresdb, queries, jwtService, cookieService, authMiddleware)

	// metrics 등록
	metrics.Register(app)

	// 기본 라우트
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(response.OK("Application is running PORT : "+cfg.App.Port, nil))
	})

	// application 실행
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Error("server stopped", log.MapErr("error", err))
	}
	log.Info("Application is running", log.MapStr("port", cfg.App.Port))
}
