package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
	"github.com/andrey67895/L0_TEST_TASK/internal/logger"
)

type RedisOrderCache struct {
	log    *logger.Logger
	client *redis.Client
}

func NewRedisOrderCache(client *redis.Client, log *logger.Logger) *RedisOrderCache {
	return &RedisOrderCache{client: client, log: log}
}

func (r *RedisOrderCache) GetByID(ctx context.Context, uid string) (domain.Order, bool) {
	val, err := r.client.Get(ctx, uid).Result()
	if err == redis.Nil {
		r.log.Error("Не найдена информация о заказе в Redis")
		return domain.Order{}, false
	} else if err != nil {
		r.log.Error("Ошибка при получении информация о заказе в Redis")
		return domain.Order{}, false
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return domain.Order{}, false
	}
	r.log.Info("Получение информации о заказе из Redis")
	return order, true
}

func (r *RedisOrderCache) SetByID(ctx context.Context, order domain.Order, ttl time.Duration) {
	data, err := json.Marshal(order)
	if err != nil {
		// логируем ошибку
		return
	}

	if err := r.client.Set(ctx, order.OrderUID, data, ttl).Err(); err != nil {
		// логируем ошибку
		return
	}
}
