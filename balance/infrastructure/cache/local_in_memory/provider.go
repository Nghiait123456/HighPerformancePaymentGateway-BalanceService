package local_in_memory

import (
	"github.com/google/wire"
	"github.com/patrickmn/go-cache"
	"time"
)

func defaultExpiration() time.Duration {
	//todo get from global config
	var time time.Duration
	return time
}

func cleanupInterval() time.Duration {
	//todo get from global config
	var time time.Duration
	return time
}

var ProviderLocalInMemory = wire.NewSet(
	NewLocalInMemory,
	cache.New(defaultExpiration(), cleanupInterval()),
	wire.Bind(new(LocalInMemoryInterface), new(*LocalInMemory)),
)
