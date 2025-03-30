package handler

import (
	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	c.JSON(200, cache.Config)
}

func GetTraffc(c *gin.Context) {
	c.JSON(200, map[string]any{
		"Label": cache.Label,
		"TX":    cache.TXBPS,
		"RX":    cache.RXBPS,
	})
}
