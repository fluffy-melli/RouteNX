package stats

import (
	"bytes"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"io"

	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/gin-gonic/gin"
)

type response struct {
	gin.ResponseWriter
	size   int
	header http.Header
}

func (r *response) Write(data []byte) (int, error) {
	r.copyHeaders()
	n, err := r.ResponseWriter.Write(data)
	r.size += n
	return n, err
}

func (r *response) WriteHeader(code int) {
	r.copyHeaders()
	r.ResponseWriter.WriteHeader(code)
}

func (r *response) copyHeaders() {
	if r.header == nil {
		r.header = make(http.Header)
		for k, vv := range r.ResponseWriter.Header() {
			for _, v := range vv {
				r.header.Add(k, v)
			}
		}
	}
}

func TX() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := &response{ResponseWriter: c.Writer}
		c.Writer = resp
		c.Next()
		var total int64

		for name, values := range resp.header {
			for _, value := range values {
				total += int64(len(fmt.Sprintf("%s: %s\r\n", name, value)))
			}
		}

		total += 2
		total += int64(resp.size)

		atomic.AddInt64(&cache.Value.TX, total)
	}
}

func RX() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		c.Request.Body = io.NopCloser(&buf)

		body, err := io.ReadAll(tee)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()

		var total int64

		total += int64(len(fmt.Sprintf("%s %s %s\r\n", c.Request.Method, c.Request.RequestURI, c.Request.Proto)))

		for name, values := range c.Request.Header {
			for _, value := range values {
				total += int64(len(fmt.Sprintf("%s: %s\r\n", name, value)))
			}
		}

		total += 2
		total += int64(len(body))

		atomic.AddInt64(&cache.Value.RX, total)
	}
}

func Traffic() {
	sleep := 30

	ticker := time.NewTicker(time.Duration(sleep) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rxBps := (atomic.SwapInt64(&cache.Value.RX, 0) * 8) / int64(sleep)
		txBps := (atomic.SwapInt64(&cache.Value.TX, 0) * 8) / int64(sleep)
		now := time.Now().UnixMilli()

		cache.Value.Lock()

		if len(cache.Value.RXBPS) >= 20 {
			cache.Value.RXBPS = cache.Value.RXBPS[1:]
		}
		cache.Value.RXBPS = append(cache.Value.RXBPS, rxBps)

		if len(cache.Value.TXBPS) >= 20 {
			cache.Value.TXBPS = cache.Value.TXBPS[1:]
		}
		cache.Value.TXBPS = append(cache.Value.TXBPS, txBps)

		if len(cache.Value.Label) >= 20 {
			cache.Value.Label = cache.Value.Label[1:]
		}
		cache.Value.Label = append(cache.Value.Label, now)

		cache.Value.Unlock()
	}
}
