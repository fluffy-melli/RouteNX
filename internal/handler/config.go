package handler

import (
	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	c.JSON(200, cache.Config)
}
