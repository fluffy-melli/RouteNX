package tcp

import (
	"fmt"
	"strings"
)

func Protocol(buffer []byte) string {
	parts := strings.Fields(string(buffer))
	if len(parts) > 0 {
		method := parts[0]
		switch method {
		case "GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH":
			return "HTTP"
		}
	}
	return "UNN"
}

func Host(buffer []byte) string {
	lines := strings.Split(string(buffer), "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Host:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Host:"))
		}
	}
	return ""
}

func Data(buffer []byte) ([]string, []string, error) {
	lines := strings.Split(string(buffer), "\r\n")
	if len(lines) < 1 {
		return nil, nil, fmt.Errorf("error reading request")
	}
	parts := strings.Fields(lines[0])
	if len(parts) < 3 {
		return nil, nil, fmt.Errorf("error reading request")
	}
	return lines, parts, nil
}
