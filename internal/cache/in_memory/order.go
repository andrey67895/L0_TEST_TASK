package in_memory

import (
	"context"
	"sync"
	"time"

	"github.com/andrey67895/L0_TEST_TASK/internal/domain"
)

type cacheItem struct {
	value     domain.Order
	expiresAt time.Time
}

type InMemoryCache struct {
	mu       sync.RWMutex
	capacity int
	items    map[string]cacheItem
	ring     []string
	head     int
	tail     int
	size     int
	ticker   *time.Ticker
	stop     chan struct{}
}

func NewInMemoryCache(capacity int, cleanupInterval time.Duration) *InMemoryCache {
	c := &InMemoryCache{
		capacity: capacity,
		items:    make(map[string]cacheItem, capacity),
		ring:     make([]string, capacity),
		ticker:   time.NewTicker(cleanupInterval),
		stop:     make(chan struct{}),
	}

	go c.cleanupLoop()
	return c
}

func (c *InMemoryCache) cleanupLoop() {
	for {
		select {
		case <-c.ticker.C:
			c.cleanup()
		case <-c.stop:
			return
		}
	}
}

func (c *InMemoryCache) cleanup() {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if now.After(item.expiresAt) {
			delete(c.items, key)
			c.size--
		}
	}
}

func (c *InMemoryCache) Stop() {
	close(c.stop)
	c.ticker.Stop()
}

func (c *InMemoryCache) SetByID(_ context.Context, value domain.Order, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	if item, ok := c.items[value.OrderUID]; ok {
		if now.After(item.expiresAt) {
			delete(c.items, value.OrderUID)
			c.size--
		} else {
			c.items[value.OrderUID] = cacheItem{
				value:     value,
				expiresAt: now.Add(ttl),
			}
			return
		}
	}

	if c.size == c.capacity {
		oldKey := c.ring[c.head]
		delete(c.items, oldKey)
		c.head = (c.head + 1) % c.capacity
		c.size--
	}

	c.items[value.OrderUID] = cacheItem{
		value:     value,
		expiresAt: now.Add(ttl),
	}
	c.ring[c.tail] = value.OrderUID
	c.tail = (c.tail + 1) % c.capacity
	c.size++
}

func (c *InMemoryCache) GetByID(_ context.Context, key string) (domain.Order, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().After(item.expiresAt) {
		return domain.Order{}, false
	}

	return item.value, true
}
