package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

type ClientResponse struct {
	Header     http.Header
	Body       []byte
	StatusCode int
}

type Client struct {
	HTTPClient *http.Client
	Request    func(context.Context, string, string, []byte) (*http.Request, error)
}

func NewClient() *Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.Logger = LogZ{log.With().Str("component", "request").Logger()}

	return &Client{
		HTTPClient: client.StandardClient(),
		Request:    NewRequest,
	}
}

func (c *Client) Send(ctx context.Context, URL, method string, headers map[string]interface{}, payload []byte) (*ClientResponse, error) {
	req, err := c.Request(ctx, URL, method, payload)
	if err != nil {
		return nil, err //nolint:wrapcheck // not need here
	}

	for k, v := range headers {
		req.Header.Add(k, fmt.Sprint(v))
	}

	req.Header.Add("Content-Length", fmt.Sprint(len(payload)))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request [%s]: %w", URL, err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}, err
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
