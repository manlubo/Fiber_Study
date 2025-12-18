package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	zap *zap.Logger
}

var instance *logger

// Init : main에서 한 번만 호출
// ex) LOG_LEVEL=DEBUG / INFO / WARN / ERROR
func Init(env string, level string) {
	lv := strings.ToLower(level)

	var zapLevel zapcore.Level
	if err := zapLevel.Set(lv); err != nil {
		zapLevel = zapcore.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:     "timestamp",
		LevelKey:    "level",
		MessageKey:  "message",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	var encoder zapcore.Encoder

	// dev → 컬러 콘솔
	if strings.ToLower(env) == "dev" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// prod → JSON
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	instance = &logger{
		zap: zap.New(core),
	}
}

// 기본값 설정
func ensure() {
	if instance == nil {
		Init("dev", "debug")
	}
}

func MapStr(key string, data string) zap.Field {
	return zap.String(key, data)
}

func MapInt64(key string, data int64) zap.Field {
	return zap.Int64(key, data)
}

func MapInt(key string, data int) zap.Field {
	return zap.Int(key, data)
}

func MapErr(key string, err error) zap.Field {
	return zap.Error(err)
}

func MapBool(key string, data bool) zap.Field {
	return zap.Bool(key, data)
}

func MapAny(key string, data any) zap.Field {
	return zap.Any(key, data)
}

// 스타일 공통 로직
func Debug(msg string, fields ...zap.Field) {
	ensure()
	instance.zap.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	ensure()
	instance.zap.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	ensure()
	instance.zap.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	ensure()
	instance.zap.Error(msg, fields...)
}
