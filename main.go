package main

import (
	"fmt"

	"github.com/fluffy-melli/RouteNX/internal/console"
	"github.com/fluffy-melli/RouteNX/internal/proxy"
	"github.com/fluffy-melli/RouteNX/internal/proxy/middleware"
	"github.com/fluffy-melli/RouteNX/pkg/cache"
)

func main() {
	cache := cache.NewCache()

	go func() {
		router := proxy.Router(cache)
		router.Run(fmt.Sprintf(":%d", cache.Config.Port))
	}()

	go func() {
		router := console.Router(cache)
		router.Run(fmt.Sprintf(":%d", cache.Config.WebPort))
	}()

	middleware.Traffic(cache)
}
