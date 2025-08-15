package interfaces

import (
	"context"
	"time"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
)

type OrderRepository interface {
	GetByID(ctx context.Context, uid string) (domain.Order, error)
	Create(ctx context.Context, order domain.Order) error
	GetLastNOrders(ctx context.Context, n int) (orders []domain.Order, err error)
}

type OrderCacheRepository interface {
	GetByID(ctx context.Context, uid string) (domain.Order, bool)
	SetByID(ctx context.Context, order domain.Order, ttl time.Duration)
}
