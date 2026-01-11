package cache

import (
	"sync"
	"time"
)

// Cache 简单的内存缓存，支持 TTL
type Cache struct {
	mu      sync.RWMutex
	data    map[string]*item
	janitor *time.Ticker
	stop    chan struct{}
}

type item struct {
	value    interface{}
	expireAt time.Time
}

// New 创建新的缓存实例
// cleanupInterval 指定清理过期项的间隔时间
func New(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		data:    make(map[string]*item),
		janitor: time.NewTicker(cleanupInterval),
		stop:    make(chan struct{}),
	}
	go c.cleanup()
	return c
}

// Get 获取缓存值
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, ok := c.data[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(it.expireAt) {
		return nil, false
	}

	return it.value, true
}

// Set 设置缓存值
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &item{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
}

// Delete 删除缓存值
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// Clear 清空所有缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*item)
}

// Stop 停止缓存清理协程
func (c *Cache) Stop() {
	c.janitor.Stop()
	close(c.stop)
}

// cleanup 定期清理过期项
func (c *Cache) cleanup() {
	for {
		select {
		case <-c.janitor.C:
			c.mu.Lock()
			now := time.Now()
			for k, v := range c.data {
				if now.After(v.expireAt) {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}
