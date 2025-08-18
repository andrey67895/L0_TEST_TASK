package in_memory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
)

func TestInMemoryCache_SetAndGet(t *testing.T) {
	ctx := t.Context()
	cache := NewInMemoryCache(2, time.Second*10)
	defer cache.Stop()

	order := domain.Order{OrderUID: "1"}

	cache.SetByID(ctx, order, time.Second*5)
	got, ok := cache.GetByID(ctx, "1")
	assert.True(t, ok)
	assert.Equal(t, order, got)

	// Проверка несуществующего ключа
	_, ok = cache.GetByID(ctx, "2")
	assert.False(t, ok)
}

func TestInMemoryCache_TTLExpiry(t *testing.T) {
	ctx := t.Context()
	cache := NewInMemoryCache(2, time.Millisecond*50)
	defer cache.Stop()

	order := domain.Order{OrderUID: "1"}
	cache.SetByID(ctx, order, time.Millisecond*100)

	time.Sleep(time.Millisecond * 150) // Ждем истечения TTL

	_, ok := cache.GetByID(ctx, "1")
	assert.False(t, ok, "элемент должен быть удален после истечения TTL")
}

func TestInMemoryCache_CapacityEviction(t *testing.T) {
	ctx := t.Context()
	cache := NewInMemoryCache(2, time.Second*10)
	defer cache.Stop()

	order1 := domain.Order{OrderUID: "1"}
	order2 := domain.Order{OrderUID: "2"}
	order3 := domain.Order{OrderUID: "3"}

	cache.SetByID(ctx, order1, time.Second*5)
	cache.SetByID(ctx, order2, time.Second*5)

	// При добавлении третьего, первый должен быть удален
	cache.SetByID(ctx, order3, time.Second*5)

	_, ok := cache.GetByID(ctx, "1")
	assert.False(t, ok, "первый элемент должен быть удален из-за переполнения")

	got, ok := cache.GetByID(ctx, "2")
	assert.True(t, ok)
	assert.Equal(t, order2, got)

	got, ok = cache.GetByID(ctx, "3")
	assert.True(t, ok)
	assert.Equal(t, order3, got)
}

func TestInMemoryCache_UpdateExisting(t *testing.T) {
	ctx := t.Context()
	cache := NewInMemoryCache(2, time.Second*10)
	defer cache.Stop()

	order := domain.Order{OrderUID: "1", Locale: "RU"}
	cache.SetByID(ctx, order, time.Second*5)

	// Обновляем тот же ключ
	updatedOrder := domain.Order{OrderUID: "1"}
	cache.SetByID(ctx, updatedOrder, time.Second*5)

	got, ok := cache.GetByID(ctx, "1")
	assert.True(t, ok)
	assert.Equal(t, updatedOrder, got)
}
