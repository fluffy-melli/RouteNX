package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"io"

	"github.com/fluffy-melli/RouteNX/pkg/cache"
	"github.com/gin-gonic/gin"
)

type response struct {
	gin.ResponseWriter
	size int
}

func (w *response) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	w.size += n
	return n, err
}

func TX(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &response{ResponseWriter: c.Writer}
		c.Writer = resp
		c.Next()
		atomic.AddInt64(&cache.TX, int64(resp.size))
	}
}

func RX(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		atomic.AddInt64(&cache.RX, int64(len(fmt.Sprintf("%s %s %s\r\n", c.Request.Method, c.Request.RequestURI, c.Request.Proto))))
		for name, values := range c.Request.Header {
			for _, value := range values {
				atomic.AddInt64(&cache.RX, int64(len(fmt.Sprintf("%s: %s\r\n", name, value))))
			}
		}
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, err := io.ReadAll(tee)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Request.Body = io.NopCloser(&buf)
		atomic.AddInt64(&cache.RX, int64(len(body))+2)
		c.Next()
	}
}

func Traffic(cache *cache.Cache) {
	sleep := 30

	ticker := time.NewTicker(time.Duration(sleep) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rxBps := (atomic.SwapInt64(&cache.RX, 0) * 8) / int64(sleep)
		txBps := (atomic.SwapInt64(&cache.TX, 0) * 8) / int64(sleep)
		now := time.Now().UnixMilli()

		cache.Lock()

		if len(cache.RXBPS) >= 20 {
			cache.RXBPS = cache.RXBPS[1:]
		}
		cache.RXBPS = append(cache.RXBPS, rxBps)

		if len(cache.TXBPS) >= 20 {
			cache.TXBPS = cache.TXBPS[1:]
		}
		cache.TXBPS = append(cache.TXBPS, txBps)

		if len(cache.Label) >= 20 {
			cache.Label = cache.Label[1:]
		}
		cache.Label = append(cache.Label, now)

		cache.Unlock()
	}
}
