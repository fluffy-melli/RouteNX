package middleware

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/fluffy-melli/RouteNX/internal/cache"
)

func Receive(buffer []byte) {
	atomic.AddInt64(&cache.RX, int64(len(buffer)))
}
func Transmit(conn net.Conn, buffer []byte) {
	atomic.AddInt64(&cache.TX, int64(len(buffer)))
	conn.Write(buffer)
}

func BPS() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		rxBps := atomic.LoadInt64(&cache.RX) * 8
		txBps := atomic.LoadInt64(&cache.TX) * 8

		cache.Label = append(cache.Label, time.Now().UnixMilli())
		cache.RXBPS = append(cache.RXBPS, rxBps)
		cache.TXBPS = append(cache.TXBPS, txBps)

		if len(cache.RXBPS) > 20 {
			cache.RXBPS = cache.RXBPS[1:]
		}

		if len(cache.TXBPS) > 20 {
			cache.TXBPS = cache.TXBPS[1:]
		}

		if len(cache.Label) > 20 {
			cache.Label = cache.Label[1:]
		}

		atomic.StoreInt64(&cache.RX, 0)
		atomic.StoreInt64(&cache.TX, 0)
	}
}
