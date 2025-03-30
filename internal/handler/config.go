package handler

import (
	"net/http"

	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/fluffy-melli/RouteNX/pkg/config"
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

func PutConfig(c *gin.Context) {
	var configs config.RouteNX
	if err := c.ShouldBindJSON(&configs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cache.Config = &configs
	go cache.Config.SaveToFile(config.RouteNXJSON)
	c.JSON(http.StatusOK, cache.Config)
}
