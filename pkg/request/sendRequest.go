package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	timeout = 5
)

type Client struct {
	HTTPClient *http.Client
	Request    func(context.Context, string, string, []byte) (*http.Request, error)
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout:   timeout * time.Second,
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		},
		Request: NewRequest,
	}
}

func (c *Client) Send(ctx context.Context, URL, method string, headers map[string]string, payload []byte) ([]byte, error) {
	req, err := c.Request(ctx, URL, method, payload)
	if err != nil {
		return nil, err //nolint:wraperr // not need here
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request [%s]: %w", URL, err)
	}

	// _, _ = io.Copy(io.Discard, resp.Body)
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	return body, err
}

func NewRequest(ctx context.Context, URL, method string, payload []byte) (*http.Request, error) {
	var body io.Reader

	if payload != nil {
		body = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for [%s]: %w", URL, err)
	}

	return req, nil
}
