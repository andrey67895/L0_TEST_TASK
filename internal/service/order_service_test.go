package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
)

// Создаем мок версии Репозитория и Кэша

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) GetByID(ctx context.Context, uid string) (domain.Order, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetLastNOrders(ctx context.Context, n int) ([]domain.Order, error) {
	args := m.Called(ctx, n)
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockOrderRepository) Create(ctx context.Context, order domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

type MockOrderCache struct {
	mock.Mock
}

func (m *MockOrderCache) GetByID(ctx context.Context, uid string) (domain.Order, bool) {
	args := m.Called(ctx, uid)
	return args.Get(0).(domain.Order), args.Bool(1)
}

func (m *MockOrderCache) SetByID(ctx context.Context, order domain.Order, ttl time.Duration) {
	m.Called(ctx, order, ttl)
}

func TestGetOrderByUID_CacheHit(t *testing.T) {
	ctx := context.Background()
	cache := new(MockOrderCache)
	repo := new(MockOrderRepository)
	svc := NewOrderService(repo, cache)

	expectedOrder := domain.Order{OrderUID: "123"}

	cache.On("GetByID", ctx, "123").Return(expectedOrder, true)

	order, err := svc.GetOrderByUID(ctx, "123")

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	cache.AssertCalled(t, "GetByID", ctx, "123")
	repo.AssertNotCalled(t, "GetByID", ctx, "123")
}

func TestGetOrderByUID_CacheMiss(t *testing.T) {
	ctx := context.Background()
	cache := new(MockOrderCache)
	repo := new(MockOrderRepository)
	svc := NewOrderService(repo, cache)

	expectedOrder := domain.Order{OrderUID: "123"}

	cache.On("GetByID", ctx, "123").Return(domain.Order{}, false)
	repo.On("GetByID", ctx, "123").Return(expectedOrder, nil)
	cache.On("SetByID", ctx, expectedOrder, mock.Anything).Return()

	order, err := svc.GetOrderByUID(ctx, "123")

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	cache.AssertCalled(t, "SetByID", ctx, expectedOrder, mock.Anything)
}

func TestAddLastOrderInCache(t *testing.T) {
	ctx := context.Background()
	cache := new(MockOrderCache)
	repo := new(MockOrderRepository)
	svc := NewOrderService(repo, cache)

	orders := []domain.Order{
		{OrderUID: "1"},
		{OrderUID: "2"},
	}

	repo.On("GetLastNOrders", ctx, 2).Return(orders, nil)
	cache.On("SetByID", ctx, orders[0], mock.Anything).Return()
	cache.On("SetByID", ctx, orders[1], mock.Anything).Return()

	err := svc.AddLastOrderInCache(ctx, 2)
	assert.NoError(t, err)
	cache.AssertNumberOfCalls(t, "SetByID", 2)
}

func TestCreateOrder_Success(t *testing.T) {
	ctx := context.Background()
	cache := new(MockOrderCache)
	repo := new(MockOrderRepository)
	svc := NewOrderService(repo, cache)

	order := domain.Order{OrderUID: "123"}

	repo.On("Create", ctx, order).Return(nil)
	cache.On("SetByID", ctx, order, mock.Anything).Return()

	err := svc.CreateOrder(ctx, order)
	assert.NoError(t, err)
	cache.AssertCalled(t, "SetByID", ctx, order, mock.Anything)
}

func TestCreateOrder_Error(t *testing.T) {
	ctx := context.Background()
	cache := new(MockOrderCache)
	repo := new(MockOrderRepository)
	svc := NewOrderService(repo, cache)

	order := domain.Order{OrderUID: "123"}
	repoErr := errors.New("repo error")

	repo.On("Create", ctx, order).Return(repoErr)

	err := svc.CreateOrder(ctx, order)
	assert.Equal(t, repoErr, err)
	cache.AssertNotCalled(t, "SetByID", ctx, order, mock.Anything)
}
