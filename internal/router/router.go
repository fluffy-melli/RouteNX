package router

import (
	"github.com/fluffy-melli/RouteNX/internal/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Static("/static", "./dist")
	r.Use(cors.Default())
	r.GET("/route", handler.GetConfig)
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})
	return r
}
