package middleware

import (
	"net"
	"strings"

	"github.com/fluffy-melli/RouteNX/internal/cache"
	"github.com/fluffy-melli/RouteNX/pkg/proxy"
	"github.com/fluffy-melli/RouteNX/pkg/status"
	"github.com/fluffy-melli/RouteNX/pkg/tcp"
)

func Listener(conn net.Conn) {
	defer conn.Close()

	req, err := tcp.Receive(conn)
	if err != nil {
		Transmit(conn, []byte(status.S500+status.ETR))
		return
	}

	Receive(req)

	switch tcp.Protocol(req) {
	case "HTTP":
		route := cache.Config.GetRoute(tcp.Host(req))
		if route == nil {
			Transmit(conn, []byte(status.S400+status.ETR+"No Route"))
			return
		}
		ips := strings.Split(conn.RemoteAddr().String(), ":")
		if cache.Config.IsBlock(route, ips[0]) {
			Transmit(conn, []byte(status.S400+status.ETR+"CIDR IP BAN"))
			return
		}
		Transmit(conn, proxy.HTTP(conn, req, route.Endpoint))
	default:
		Transmit(conn, []byte(status.S400+status.ETR))
	}
}
