package middleware

import (
	"io"
	"net/http"

	"github.com/fluffy-melli/RouteNX/pkg/cache"
	"github.com/fluffy-melli/RouteNX/pkg/firewall"
	"github.com/fluffy-melli/RouteNX/pkg/logger"
	"github.com/fluffy-melli/RouteNX/pkg/request"
	"github.com/gin-gonic/gin"
)

func Proxy(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		to := cache.Config.GetRoute(c.Request.Host)
		if to == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Route"})
			return
		}

		if firewall.IsCidrBlock(cache.Config, to, c.RemoteIP()) {
			c.JSON(http.StatusForbidden, gin.H{"error": "CIDR IP Block"})
			return
		}

		req, err := request.HTTP(c, to.Endpoint)
		if err != nil {
			logger.WARNING("%s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.WARNING("%s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		c.Status(resp.StatusCode)

		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			logger.WARNING("%s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}
