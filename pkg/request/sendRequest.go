package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog"
)

type ClientResponse struct {
	Header     http.Header
	Body       []byte
	StatusCode int
}

type Client struct {
	HTTPClient *http.Client
}

type Config struct {
	SkipVerify bool
	Pooled     bool
	Log        *zerolog.Logger
}

var (
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

func NewClient(cfg Config) *Client {
	var logger interface{}
	if cfg.Log != nil {
		logger = LogZ{Log: *cfg.Log}
	}

	var httpClient *http.Client
	if cfg.Pooled {
		httpClient = cleanhttp.DefaultPooledClient()
	} else {
		httpClient = cleanhttp.DefaultClient()
	}

	if cfg.SkipVerify {
		//nolint:forcetypeassert // clear
		tlsClientConfig := httpClient.Transport.(*http.Transport).TLSClientConfig
		if tlsClientConfig == nil {
			tlsClientConfig = &tls.Config{
				//nolint:gosec // user defined
				InsecureSkipVerify: true,
			}
		} else {
			tlsClientConfig.InsecureSkipVerify = true
		}
		//nolint:forcetypeassert // clear
		httpClient.Transport.(*http.Transport).TLSClientConfig = tlsClientConfig
	}

	client := &retryablehttp.Client{
		HTTPClient:   httpClient,
		Logger:       logger,
		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,
		CheckRetry:   RetryPolicy,
		Backoff:      retryablehttp.DefaultBackoff,
	}

	return &Client{
		HTTPClient: client.StandardClient(),
	}
}

func (c *Client) Send(
	ctx context.Context,
	url, method string,
	headers map[string]interface{},
	payload []byte,
	retry *Retry,
	skipVerify bool,
) (*ClientResponse, error) {
	var resp *http.Response

	if retry != nil {
		ctx = context.WithValue(ctx, RetryCodesValue, retry)
	}

	req, err := c.NewRequest(ctx, url, method, headers, payload)
	if err != nil {
		return nil, err
	}

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return &ClientResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
	}, err
}

// NewRequest creates a new HTTP request with the given method, URL, and optional body.
//
//nolint:lll // clear
func (c *Client) NewRequest(ctx context.Context, url, method string, headers map[string]interface{}, payload []byte) (*http.Request, error) {
	var body io.Reader

	if payload != nil {
		body = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for [%s]: %w", url, err)
	}

	for k, v := range headers {
		req.Header.Add(k, fmt.Sprint(v))
	}

	req.Header.Add("Content-Length", fmt.Sprint(len(payload)))

	return req, nil
}
