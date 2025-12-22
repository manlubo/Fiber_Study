package config

import (
	"fmt"
	"os"
	"study/pkg/log"
	"study/pkg/util"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App           App
	Postgres      Postgres `yaml:"postgres"`
	JWT           JWT
	Log           Log           `yaml:"log"`
	Cors          Cors          `yaml:"cors"`
	Cookie        Cookie        `yaml:"cookie"`
	Observability Observability `yaml:"observability"`
}

func Load() (*Config, error) {
	var cfg Config

	// env 설정 파일 로드
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	// yaml 설정 파일 로드
	ymlPath := fmt.Sprintf(util.GetPath("configs/application-%s.yml"), cfg.App.Env)
	data, err := os.ReadFile(ymlPath)
	if err != nil {
		fmt.Println("yml 파일을 읽는데 실패했습니다 :", ymlPath)
		return nil, err
	}

	// yaml 설정 로드
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Println("yml 설정을 읽는데 실패했습니다 :", ymlPath)
		return nil, err
	}

	// 로깅 초기화
	log.Info("실행환경", log.MapStr("env", cfg.App.Env), log.MapStr("logLevel", cfg.Log.Level))
	log.Init(cfg.App.Env, cfg.Log.Level)
	log.Info("JWT", log.MapInt("accessExpireMin", cfg.JWT.AccessExpireMin), log.MapInt("refreshExpireDay", cfg.JWT.RefreshExpireDay))

	return &cfg, nil
}
