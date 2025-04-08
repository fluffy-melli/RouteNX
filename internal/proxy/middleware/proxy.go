package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/fluffy-melli/RouteNX/internal/proxy/handler"
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
			handler.NoRoute(c)
			return
		}

		if firewall.IsCidrBlock(cache.Config, to, c.RemoteIP()) {
			cache.Logger.AddBlockLog(logger.BlockLogger{
				OriginIP:  c.ClientIP(),
				ForwardIP: c.RemoteIP(),
				Host:      c.Request.Host,
				Time:      time.Now().Format(time.RFC3339),
			})
			handler.IPBlock(c)
			return
		}

		req, err := request.HTTP(c, to.Endpoint)
		if err != nil {
			cache.Logger.AddErrorLog(logger.ErrorLogger{
				Error: err.Error(),
				Time:  time.Now().Format(time.RFC3339),
			})
			logger.WARNING("%s", err.Error())
			handler.InternalError(c, err)
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			cache.Logger.AddErrorLog(logger.ErrorLogger{
				Error: err.Error(),
				Time:  time.Now().Format(time.RFC3339),
			})
			logger.WARNING("%s", err.Error())
			handler.InternalError(c, err)
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
			cache.Logger.AddErrorLog(logger.ErrorLogger{
				Error: err.Error(),
				Time:  time.Now().Format(time.RFC3339),
			})
			logger.WARNING("%s", err.Error())
			handler.InternalError(c, err)
		}
	}
}

func SSLRedirect(c *gin.Context) {
	if c.Request.TLS == nil {
		target := "https://" + c.Request.Host + c.Request.URL.String()
		c.Redirect(http.StatusMovedPermanently, target)
		c.Abort()
		return
	}
	c.Next()
}
