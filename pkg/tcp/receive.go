package tcp

import (
	"io"
	"net"
)

func Receive(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buffer[:n], nil
}
