package cache

import (
	"sync"
	"sync/atomic"

	"github.com/fluffy-melli/RouteNX/pkg/config"
)

type Cache struct {
	sync.Mutex
	Config *config.RouteNX
	Label  []int64
	RXBPS  []int64
	TXBPS  []int64
	RX     int64
	TX     int64
}

func NewCache() *Cache {
	var err error
	cfg, err := config.LoadFromFile(config.RouteNXJSON)
	if err != nil {
		cfg = config.NewRouteNX()
		cfg.SaveToFile(config.RouteNXJSON)
	}

	return &Cache{
		Config: cfg,
	}
}

func (c *Cache) AddRX(value int64) {
	atomic.AddInt64(&c.RX, value)
}

func (c *Cache) AddTX(value int64) {
	atomic.AddInt64(&c.TX, value)
}
