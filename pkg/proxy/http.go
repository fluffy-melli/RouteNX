package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/fluffy-melli/RouteNX/pkg/request"
	"github.com/fluffy-melli/RouteNX/pkg/status"
	"github.com/fluffy-melli/RouteNX/pkg/tcp"
)

func HTTP(conn net.Conn, buffer []byte, target string) []byte {
	_, urls, err := tcp.Data(buffer)
	if err != nil {
		return []byte(status.S400 + status.ETR)
	}

	resp, err := request.HTTP(urls[0], fmt.Sprintf("%s%s", target, urls[1]), nil, nil)
	if err != nil {
		return []byte(status.S500 + status.ETR)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s %d %s\r\n", resp.Proto, resp.StatusCode, http.StatusText(resp.StatusCode))

	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Fprintf(&buf, "%s: %s\r\n", key, value)
		}
	}

	buf.WriteString("\r\n")

	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return []byte(status.S500 + status.ETR)
	}

	return buf.Bytes()
}
