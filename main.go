package main

import (
	"fmt"
	"net"
	"time"

	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/fluffy-melli/RouteNX/internal/console"
	"github.com/fluffy-melli/RouteNX/internal/proxy"
	"github.com/fluffy-melli/RouteNX/pkg/certs"
	"github.com/fluffy-melli/RouteNX/pkg/logs"
	"github.com/fluffy-melli/RouteNX/pkg/stats"
	"github.com/gin-gonic/gin"
)

func Run(router *gin.Engine, port uint16, server string) {
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		logs.INFO("%s failed to start: %s", server, err.Error())
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
	cache.Value = cache.NewCache()

	logs.TRY("--------------------- [{yellow}URL{reset}] ----------------------")
	logs.INFO("proxy (http)  : {blue}http://%s:%d{reset}", LocalIP(), cache.Value.Config.Port)
	logs.INFO("proxy (https) : {blue}http://%s:%d{reset}", LocalIP(), cache.Value.Config.SSLPort)
	logs.INFO("web-console   : {blue}http://%s:%d{reset}", LocalIP(), cache.Value.Config.WebPort)

	go func() {
		router := proxy.Router()
		if cache.Value.Config.SSL.Enabled {
			logs.TRY("--------------- [{yellow}Trying: SSL Cert{reset}] ---------------")
			SSL, err := certs.NewSSL(cache.Value.Config.SSL.Domains, cache.Value.Config.SSL.Email, cache.Value.Config.SSL.Testing)
			if err != nil {
				logs.ERROR("--------------- [{red}Failed: SSL Cert{reset}] ---------------\n{red}%v{reset}", err)
				return
			}
			go func(SSL *certs.SSL) {
				for {
					time.Sleep(60 * 24 * time.Hour)
					logs.TRY("--------------- [{yellow}Trying: SSL Renew{reset}] --------------")
					if err := SSL.Renew(); err != nil {
						logs.ERROR("--------------- [{red}Failed: SSL Renew{reset}] --------------\n{red}%v{reset}", err)
					}
				}
			}(SSL)
			if err := SSL.ApplyToGin(router, fmt.Sprintf(":%d", cache.Value.Config.SSLPort)); err != nil {
				logs.ERROR("Proxy server (ssl) failed to start: %s", err.Error())
				return
			}
		}
		if err := router.Run(fmt.Sprintf(":%d", cache.Value.Config.Port)); err != nil {
			logs.ERROR("Proxy server failed to start: %s", err.Error())
			return
		}
	}()

	go func() {
		router := console.Router()
		if err := router.Run(fmt.Sprintf(":%d", cache.Value.Config.WebPort)); err != nil {
			logs.ERROR("Web console server failed to start: %s", err.Error())
			return
		}
	}()

	stats.Traffic()
}
