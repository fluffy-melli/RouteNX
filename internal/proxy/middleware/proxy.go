package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/fluffy-melli/RouteNX/internal/proxy/handler"
	"github.com/fluffy-melli/RouteNX/pkg/firewall"
	"github.com/fluffy-melli/RouteNX/pkg/logs"
	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context) {
	logs.CsvAppend([]string{
		time.Now().String(),
		c.RemoteIP(),
		c.ClientIP(),
		c.Request.Host,
		c.Request.RequestURI,
	})

	to := cache.Value.Config.GetRoute(c.Request.Host)
	if to == nil {
		handler.NoRoute(c)
		return
	}

	if firewall.IsCidrBlock(cache.Value.Config, to, c.RemoteIP()) {
		cache.Value.Logger.AddBlockLog(cache.Block{
			OriginIP:  c.ClientIP(),
			ForwardIP: c.RemoteIP(),
			Host:      c.Request.Host,
			Time:      time.Now().Format(time.RFC3339),
		})
		handler.IPBlock(c)
		return
	}

	req, err := http.NewRequest(c.Request.Method, to.Endpoint+c.Request.RequestURI, c.Request.Body)
	if err != nil {
		handler.InternalError(c, err)
		return
	}

	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.ContentLength = c.Request.ContentLength

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
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
		handler.InternalError(c, err)
		return
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
