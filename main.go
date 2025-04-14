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

	logger.TRY("--------------------- [{yellow}URL{reset}] ----------------------")
	logger.INFO("proxy (http)  : {blue}http://%s:%d{reset}", LocalIP(), cache.Config.Port)
	logger.INFO("proxy (https) : {blue}http://%s:%d{reset}", LocalIP(), cache.Config.SSLPort)
	logger.INFO("web-console   : {blue}http://%s:%d{reset}", LocalIP(), cache.Config.WebPort)

	go func() {
		router := proxy.Router(cache)
		if cache.Config.SSL.Enabled {
			logger.TRY("--------------- [{yellow}Trying: SSL Cert{reset}] ---------------")
			SSL, err := ssl.NewSSL(cache.Config.SSL.Domains, cache.Config.SSL.Email, cache.Config.SSL.Testing)
			if err != nil {
				logger.ERROR("--------------- [{red}Failed: SSL Cert{reset}] ---------------\n{red}%v{reset}", err)
				return
			}
			go func(SSL *ssl.SSL) {
				for {
					time.Sleep(60 * 24 * time.Hour)
					logger.TRY("--------------- [{yellow}Trying: SSL Renew{reset}] --------------")
					if err := SSL.Renew(); err != nil {
						logger.ERROR("--------------- [{red}Failed: SSL Renew{reset}] --------------\n{red}%v{reset}", err)
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
		if err := router.Run(fmt.Sprintf(":%d", cache.Config.WebPort)); err != nil {
			logger.ERROR("Web console server failed to start: %s", err.Error())
			return
		}
	}()

	middleware.Traffic(cache)
}
