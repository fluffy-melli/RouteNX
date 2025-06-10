package console

import (
	"github.com/fluffy-melli/RouteNX/internal/console/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(cors.Default())
	r.GET("/config", handler.GetConfig)
	r.PUT("/config", handler.PutConfig)
	r.GET("/traffic", handler.GetTraffic)
	r.GET("/logger", handler.GetLogger)
	r.Static("/static", "./dist")
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})
	return r
}
