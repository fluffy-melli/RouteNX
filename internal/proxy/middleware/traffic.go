package middleware

import (
	"bytes"
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

func MeasureTraffic(cache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqSize := c.Request.ContentLength
		if reqSize < 0 {
			body, _ := io.ReadAll(c.Request.Body)
			reqSize = int64(len(body))
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
		}
		atomic.AddInt64(&cache.RX, reqSize)

		rw := &response{ResponseWriter: c.Writer}
		c.Writer = rw
		c.Next()

		atomic.AddInt64(&cache.TX, int64(rw.size))
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
