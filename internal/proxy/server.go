package proxy

import (
	"github.com/fluffy-melli/RouteNX/internal/proxy/middleware"
	"github.com/fluffy-melli/RouteNX/pkg/cache"
	"github.com/gin-gonic/gin"
)

func Router(cache *cache.Cache) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if cache.Config.SSL.Enabled {
		r.Use(middleware.SSLRedirect)
	}
	r.Use(middleware.RX(cache))
	r.Use(middleware.TX(cache))
	r.Any("/*all", middleware.Proxy(cache))
	return r
}
