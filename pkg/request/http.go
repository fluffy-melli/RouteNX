package request

import (
	"io"
	"net/http"
)

func HTTP(method, targetURL string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
