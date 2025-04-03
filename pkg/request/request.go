package request

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HTTP(c *gin.Context, endpoint string) (*http.Request, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(c.Request.Method, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	target, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	req.Host = target.Host
	req.URL.Host = target.Host
	req.URL.Path = c.Param("all")
	req.URL.Scheme = target.Scheme
	req.Header.Set("X-Real-IP", c.ClientIP())

	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return req, nil
}
