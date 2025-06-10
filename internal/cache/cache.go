package cache

import (
	"sync"
	"sync/atomic"

	"github.com/fluffy-melli/RouteNX/pkg/config"
	"github.com/fluffy-melli/RouteNX/pkg/logs"
)

type Cache struct {
	sync.Mutex
	Logger *Logger
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
		logs.WARNING(err.Error())
		cfg = config.NewRouteNX()
		cfg.SaveToFile(config.RouteNXJSON)
	}

	return &Cache{
		Logger: NewLogger(),
		Config: cfg,
	}
}

func (c *Cache) AddRX(value int64) {
	atomic.AddInt64(&c.RX, value)
}

func (c *Cache) AddTX(value int64) {
	atomic.AddInt64(&c.TX, value)
}

var Value *Cache
