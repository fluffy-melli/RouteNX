package main

import (
	"fmt"
	"net"
	"time"

	"github.com/fluffy-melli/RouteNX/internal/console"
	"github.com/fluffy-melli/RouteNX/internal/proxy"
	"github.com/fluffy-melli/RouteNX/internal/proxy/middleware"
	"github.com/fluffy-melli/RouteNX/internal/ssl"
	"github.com/fluffy-melli/RouteNX/pkg/cache"
	"github.com/fluffy-melli/RouteNX/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Run(router *gin.Engine, port uint16, server string) {
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.INFO("%s failed to start: %s", server, err.Error())
	}
}

func LocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}

func main() {
	cache := cache.NewCache()

	go func() {
		router := proxy.Router(cache)
		logger.INFO("Proxy server running at {blue}http://%s:%d{reset}", LocalIP(), cache.Config.Port)
		if cache.Config.SSL.Enabled {
			SSL, err := ssl.NewSSL(cache.Config.SSL.Domains, cache.Config.SSL.Email)
			if err != nil {
				logger.ERROR("Failed to create SSL certificate: %s", err.Error())
				return
			}
			logger.INFO("Proxy server (ssl) running at {blue}https://%s:%d{reset}", LocalIP(), cache.Config.SSLPort)
			go func(SSL *ssl.SSL) {
				for {
					time.Sleep(60 * 24 * time.Hour)
					if err := SSL.Renew(); err != nil {
						logger.ERROR("Failed to renew SSL certificate: %s", err.Error())
					} else {
						logger.INFO("SSL certificate successfully renewed")
					}
				}
			}(SSL)
			if err := SSL.ApplyToGin(router, fmt.Sprintf(":%d", cache.Config.SSLPort)); err != nil {
				logger.ERROR("Proxy server (ssl) failed to start: %s", err.Error())
				return
			}
		}
		if err := router.Run(fmt.Sprintf(":%d", cache.Config.Port)); err != nil {
			logger.ERROR("Proxy server failed to start: %s", err.Error())
			return
		}
	}()

	go func() {
		router := console.Router(cache)
		logger.INFO("Web console server running at {blue}http://%s:%d{reset}", LocalIP(), cache.Config.WebPort)
		if err := router.Run(fmt.Sprintf(":%d", cache.Config.WebPort)); err != nil {
			logger.ERROR("Web console server failed to start: %s", err.Error())
			return
		}
	}()

	middleware.Traffic(cache)
}
