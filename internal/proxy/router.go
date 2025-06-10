package proxy

import (
	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/fluffy-melli/RouteNX/internal/proxy/middleware"
	"github.com/fluffy-melli/RouteNX/pkg/stats"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if cache.Value.Config.SSL.Enabled {
		r.Use(middleware.SSLRedirect)
	}
	r.LoadHTMLGlob("templates/*")
	r.Use(stats.RX())
	r.Use(stats.TX())
	r.Any("/*all", middleware.Proxy)
	return r
}
