package logger

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// Тест создания логгера с валидной конфигурацией
func TestNewLogger_ValidConfig(t *testing.T) {
	cfg := Config{
		Level:        "info",
		Format:       "json",
		ServiceName:  "my-service",
		Environment:  "production",
		EnableCaller: true,
	}

	log, err := New(cfg)
	assert.NoError(t, err)
	assert.Equal(t, zapcore.InfoLevel, log.Level())
}

// Тест уровня логирования при неверной конфигурации
func TestNewLogger_InvalidLevel(t *testing.T) {
	cfg := Config{
		Level:        "invalid",
		Format:       "json",
		ServiceName:  "svc",
		Environment:  "production",
		EnableCaller: false,
	}

	_, err := New(cfg)
	assert.Error(t, err)
}

// Тест формата для development среды
func TestNewLogger_DevelopmentFormat(t *testing.T) {
	cfg := Config{
		Level:        "debug",
		Format:       "console",
		ServiceName:  "dev-service",
		Environment:  "development",
		EnableCaller: false,
	}

	log, err := New(cfg)
	assert.NoError(t, err)
	assert.Equal(t, zapcore.DebugLevel, log.Level())
}

// Тест функции getHostname
func TestGetHostname(t *testing.T) {
	hostname := getHostname()
	assert.NotEmpty(t, hostname)

	// Мок ошибки os.Hostname
	old := osHostname
	defer func() { osHostname = old }()
	osHostname = func() (string, error) { return "", errors.New("fail") }

	h := getHostname()
	assert.Equal(t, "unknown", h)
}

// Тест фактического логирования через observer
func TestLogger_Observer(t *testing.T) {
	core, observed := observer.New(zapcore.InfoLevel)
	log := &Logger{Logger: zap.New(core)}

	log.Info("hello world")

	assert.Equal(t, 1, observed.Len())

	entry := observed.All()[0]
	assert.Equal(t, "hello world", entry.Message)
	assert.Equal(t, zapcore.InfoLevel, entry.Level)
}
