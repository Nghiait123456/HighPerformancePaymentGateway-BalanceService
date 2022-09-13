package local_in_memory

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type LocalInMemory struct {
	cache *cache.Cache
}

type LocalInMemoryInterface interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string) (interface{}, bool)
	Init(defaultExpiration, cleanupInterval time.Duration)
}

func (l *LocalInMemory) Set(k string, x interface{}, d time.Duration) {
	l.cache.Set(k, x, d)
}

func (l *LocalInMemory) Get(k string) (interface{}, bool) {
	return l.cache.Get(k)
}

func (l *LocalInMemory) Init(defaultExpiration, cleanupInterval time.Duration) {
	l.cache = cache.New(defaultExpiration, cleanupInterval)
}

func NewLocalInMemory(cache *cache.Cache) LocalInMemoryInterface {
	return &LocalInMemory{
		cache: cache,
	}
}
