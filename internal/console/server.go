package console

import (
	"github.com/fluffy-melli/RouteNX/internal/console/handler"
	"github.com/fluffy-melli/RouteNX/pkg/cache"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router(cache *cache.Cache) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(cors.Default())
	r.GET("/config", handler.GetConfig(cache))
	r.PUT("/config", handler.PutConfig(cache))
	r.GET("/traffic", handler.GetTraffic(cache))
	r.GET("/logger", handler.GetLogger(cache))
	r.Static("/static", "./dist")
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})
	return r
}
