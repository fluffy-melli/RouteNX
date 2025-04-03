package handler

import (
	"net/http"

	"github.com/fluffy-melli/RouteNX/pkg/cache"
	"github.com/fluffy-melli/RouteNX/pkg/config"
	"github.com/gin-gonic/gin"
)

func GetConfig(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, cache.Config)
	}
}

func GetTraffic(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		cache.Lock()
		defer cache.Unlock()

		c.JSON(http.StatusOK, map[string]any{
			"Label": cache.Label,
			"TX":    cache.TXBPS,
			"RX":    cache.RXBPS,
		})
	}
}

func PutConfig(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		var configs config.RouteNX
		if err := c.ShouldBindJSON(&configs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cache.Lock()
		cache.Config = &configs
		cache.Unlock()

		go cache.Config.SaveToFile(config.RouteNXJSON)
		c.JSON(http.StatusOK, cache.Config)
	}
}
