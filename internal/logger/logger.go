package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

type Config struct {
	Level        string
	Format       string
	ServiceName  string
	Environment  string
	EnableCaller bool
}

func New(cfg Config) (*Logger, error) {
	var zapCfg zap.Config

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configure based on environment and format
	if cfg.Environment == "production" || cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
		zapCfg.EncoderConfig = encoderConfig
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig = encoderConfig
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		zapCfg.EncoderConfig.ConsoleSeparator = " | "
		zapCfg.Encoding = "console"
	}

	logLevel, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	zapCfg.Level = zap.NewAtomicLevelAt(logLevel)
	zapCfg.Development = cfg.Environment != "production"
	zapCfg.DisableCaller = !cfg.EnableCaller

	zapLogger, err := zapCfg.Build(
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("service", cfg.ServiceName),
			zap.String("environment", cfg.Environment),
			zap.String("hostname", getHostname()),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: zapLogger,
	}, nil
}

var osHostname = os.Hostname

func getHostname() string {
	hostname, err := osHostname() // <- вместо os.Hostname() теперь вызываем переменную
	if err != nil {
		return "unknown"
	}
	return hostname
}
