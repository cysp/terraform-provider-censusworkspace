package provider

import (
	"net/http"
)

type HttpClientWithUserAgent struct {
	client *http.Client

	UserAgent string
}

func NewHttpClientWithUserAgent(client *http.Client, userAgent string) *HttpClientWithUserAgent {
	return &HttpClientWithUserAgent{
		client:    client,
		UserAgent: userAgent,
	}
}

func (c *HttpClientWithUserAgent) Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" && c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	//nolint:wrapcheck
	return c.client.Do(req)
}
