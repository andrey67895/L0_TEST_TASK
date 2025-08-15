package service

import (
	"context"
	"time"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
	"github.com/andrey67895/L0_TEST_TASK/internal/interfaces"
)

type OrderService struct {
	orderRepository interfaces.OrderRepository
	orderCache      interfaces.OrderCacheRepository
}

func NewOrderService(orderRepository interfaces.OrderRepository, orderCache interfaces.OrderCacheRepository) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		orderCache:      orderCache,
	}
}

func (o *OrderService) GetOrderByUID(ctx context.Context, uid string) (domain.Order, error) {
	if order, ok := o.orderCache.GetByID(ctx, uid); ok {
		return order, nil
	}
	order, err := o.orderRepository.GetByID(ctx, uid)
	if err != nil {
		return domain.Order{}, err
	}
	o.orderCache.SetByID(ctx, order, time.Minute*5)
	return order, nil
}

func (o *OrderService) AddLastOrderInCache(ctx context.Context, capacity int) error {
	orders, err := o.orderRepository.GetLastNOrders(ctx, capacity)
	if err != nil {
		return err
	}
	for _, order := range orders {
		o.orderCache.SetByID(ctx, order, time.Minute*5)
	}
	return nil
}

func (o *OrderService) CreateOrder(ctx context.Context, order domain.Order) error {
	if err := o.orderRepository.Create(ctx, order); err == nil {
		o.orderCache.SetByID(ctx, order, time.Minute*5)
		return nil
	} else {
		return err
	}
}
