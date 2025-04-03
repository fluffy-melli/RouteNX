package main

import (
	"fmt"

	"github.com/fluffy-melli/RouteNX/internal/console"
	"github.com/fluffy-melli/RouteNX/internal/proxy"
	"github.com/fluffy-melli/RouteNX/internal/proxy/middleware"
	"github.com/fluffy-melli/RouteNX/pkg/cache"
	"github.com/fluffy-melli/RouteNX/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Run(router *gin.Engine, port uint16, server string) {
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.INFO("%s failed to start: %s", server, err.Error())
	}
}

func main() {
	cache := cache.NewCache()

	go func() {
		router := proxy.Router(cache)
		logger.INFO("Proxy server running at {blue}http://localhost:%d{reset}", cache.Config.Port)
		if err := router.Run(fmt.Sprintf(":%d", cache.Config.Port)); err != nil {
			logger.ERROR("Proxy server failed to start: %s", err.Error())
			return
		}
	}()

	go func() {
		router := console.Router(cache)
		logger.INFO("Web console server running at {blue}http://localhost:%d{reset}", cache.Config.WebPort)
		if err := router.Run(fmt.Sprintf(":%d", cache.Config.WebPort)); err != nil {
			logger.ERROR("Web console server failed to start: %s", err.Error())
			return
		}
	}()

	middleware.Traffic(cache)
}
