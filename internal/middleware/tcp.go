package middleware

import (
	"net"
	"strings"

	"github.com/fluffy-melli/RouteNX/pkg/config"
	"github.com/fluffy-melli/RouteNX/pkg/proxy"
	"github.com/fluffy-melli/RouteNX/pkg/status"
	"github.com/fluffy-melli/RouteNX/pkg/tcp"
)

func Listener(config *config.RouteNX, conn net.Conn) {
	defer conn.Close()

	req, err := tcp.Receive(conn)
	if err != nil {
		conn.Write([]byte(status.S500 + status.ETR))
		return
	}

	switch tcp.Protocol(req) {
	case "HTTP":
		route := config.GetRoute(tcp.Host(req))
		if route == nil {
			conn.Write([]byte(status.S400 + status.ETR + "No Route"))
			return
		}
		ips := strings.Split(conn.RemoteAddr().String(), ":")
		if config.IsBlock(route, ips[0]) {
			conn.Write([]byte(status.S400 + status.ETR + "CIDR IP BAN"))
			return
		}
		proxy.HTTP(conn, req, route.Endpoint)
	default:
		conn.Write([]byte(status.S400 + status.ETR))
	}
}
