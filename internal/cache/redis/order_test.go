package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
)

func TestRedisOrderCache_GetByID(t *testing.T) {
	ctx := context.Background()

	db, mock := redismock.NewClientMock()
	cache := NewRedisOrderCache(db, getLog())

	order := domain.Order{OrderUID: "123"}
	data, _ := json.Marshal(order)

	t.Run("Заказ существует", func(t *testing.T) {
		mock.ExpectGet("123").SetVal(string(data))

		got, ok := cache.GetByID(ctx, "123")
		assert.True(t, ok)
		assert.Equal(t, order, got)
	})

	t.Run("Заказ не найден", func(t *testing.T) {
		mock.ExpectGet("456").RedisNil()

		got, ok := cache.GetByID(ctx, "456")
		assert.False(t, ok)
		assert.Equal(t, domain.Order{}, got)
	})

	t.Run("Ошибка redis", func(t *testing.T) {
		mock.ExpectGet("789").SetErr(errors.New("some redis error"))

		got, ok := cache.GetByID(ctx, "789")
		assert.False(t, ok)
		assert.Equal(t, domain.Order{}, got)
	})

	t.Run("Ошибка json", func(t *testing.T) {
		mock.ExpectGet("000").SetVal("invalid json")

		got, ok := cache.GetByID(ctx, "000")
		assert.False(t, ok)
		assert.Equal(t, domain.Order{}, got)
	})
}

func TestRedisOrderCache_SetByID(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	cache := NewRedisOrderCache(db, getLog())

	order := domain.Order{OrderUID: "123"}
	data, _ := json.Marshal(order)
	ttl := time.Minute

	t.Run("Записываем успешно", func(t *testing.T) {
		mock.ExpectSet(order.OrderUID, data, ttl).SetVal("OK")

		cache.SetByID(ctx, order, ttl)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Записываем с ошибкой", func(t *testing.T) {
		mock.ExpectSet(order.OrderUID, data, ttl).SetErr(errors.New("set error"))

		cache.SetByID(ctx, order, ttl)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func getLog() *logger.Logger {
	log, _ := logger.New(logger.Config{
		Level:        "debug",
		Format:       "console",
		ServiceName:  "L0_TASK_TEST",
		Environment:  "development",
		EnableCaller: true,
	})
	return log
}
