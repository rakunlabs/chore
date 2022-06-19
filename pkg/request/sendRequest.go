package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
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
	Request    func(context.Context, string, string, map[string]interface{}, []byte) (*http.Request, error)
}

var (
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

func NewClient() *Client {
	client := &retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultClient(),
		Logger:       LogZ{log.With().Str("component", "request").Logger()},
		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,
		CheckRetry:   RetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}

	return &Client{
		HTTPClient: client.StandardClient(),
		Request:    NewRequest,
	}
}

func (c *Client) Send(
	ctx context.Context,
	URL, method string,
	headers map[string]interface{},
	payload []byte,
	retry *Retry,
) (*ClientResponse, error) {
	var resp *http.Response

	if retry != nil {
		ctx = context.WithValue(ctx, RetryCodesValue, retry)
	}

	req, err := c.Request(ctx, URL, method, headers, payload)
	if err != nil {
		return nil, err //nolint:wrapcheck // not need here
	}

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request [%s]: %w", URL, err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}, err
}

func NewRequest(ctx context.Context, URL, method string, headers map[string]interface{}, payload []byte) (*http.Request, error) {
	var body io.Reader

	if payload != nil {
		body = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, URL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for [%s]: %w", URL, err)
	}

	for k, v := range headers {
		req.Header.Add(k, fmt.Sprint(v))
	}

	req.Header.Add("Content-Length", fmt.Sprint(len(payload)))

	return req, nil
}
