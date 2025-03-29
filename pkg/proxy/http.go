package proxy

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/fluffy-melli/RouteNX/pkg/request"
	"github.com/fluffy-melli/RouteNX/pkg/status"
	"github.com/fluffy-melli/RouteNX/pkg/tcp"
)

func HTTP(conn net.Conn, buffer []byte, target string) {
	_, urls, err := tcp.Data(buffer)
	if err != nil {
		conn.Write([]byte(status.S400 + status.ETR))
		return
	}

	resp, err := request.HTTP(urls[0], fmt.Sprintf("%s%s", target, urls[1]), nil, nil)
	if err != nil {
		conn.Write([]byte(status.S500 + status.ETR))
		return
	}
	defer resp.Body.Close()

	conn.Write(fmt.Appendf([]byte{}, "%s %d %s\r\n", urls[2], resp.StatusCode, http.StatusText(resp.StatusCode)))

	for key, values := range resp.Header {
		for _, value := range values {
			conn.Write(fmt.Appendf([]byte{}, "%s: %s\r\n", key, value))
		}
	}

	conn.Write([]byte("\r\n"))

	_, err = io.Copy(conn, resp.Body)
	if err != nil {
		conn.Write([]byte(status.S500 + status.ETR))
	}
}
