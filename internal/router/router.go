package router

import (
	"github.com/fluffy-melli/RouteNX/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(cors.Default())
	r.GET("/route", handler.GetConfig)
	r.GET("/traffic", handler.GetTraffc)
	r.PUT("/config", handler.PutConfig)
	r.Static("/static", "./dist")
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})
	return r
}
