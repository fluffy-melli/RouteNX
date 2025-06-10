package handler

import (
	"net/http"

	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/fluffy-melli/RouteNX/pkg/config"
	"github.com/fluffy-melli/RouteNX/pkg/logs"
	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, cache.Value.Config)
}

func GetTraffic(c *gin.Context) {
	cache.Value.Lock()
	defer cache.Value.Unlock()

	c.JSON(http.StatusOK, map[string]any{
		"Label": cache.Value.Label,
		"TX":    cache.Value.TXBPS,
		"RX":    cache.Value.RXBPS,
	})
}

func GetLogger(c *gin.Context) {
	cache.Value.Lock()
	defer cache.Value.Unlock()

	c.JSON(http.StatusOK, cache.Value.Logger)
}

func PutConfig(c *gin.Context) {
	var configs config.RouteNX
	if err := c.ShouldBindJSON(&configs); err != nil {
		logs.WARNING("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cache.Value.Lock()
	cache.Value.Config = &configs
	cache.Value.Unlock()

	go cache.Value.Config.SaveToFile(config.RouteNXJSON)
	c.JSON(http.StatusOK, cache.Value.Config)
}
