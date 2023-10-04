package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/worldline-go/klient"
	"github.com/worldline-go/logz"
)

type Retry struct {
	Enabled             bool
	DisabledStatusCodes []int
	EnabledStatusCodes  []int
}

type ClientResponse struct {
	Header     http.Header
	Body       []byte
	StatusCode int
}

type Client struct {
	klient *klient.Client
}

type Config struct {
	SkipVerify bool
	Log        *zerolog.Logger
	Retry      Retry
	Auth       AuthConfig
}

type AuthConfig struct {
	Enabled      bool     `json:"enabled"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	TokenURL     string   `json:"token_url"`
	Scopes       []string `json:"scopes"`
}

func NewClient(cfg Config) (*Client, error) {
	options := []klient.OptionClientFn{klient.WithDisableBaseURLCheck(true)}
	if cfg.Log != nil {
		options = append(options, klient.WithLogger(logz.AdapterKV{Log: *cfg.Log}))
	}

	if cfg.SkipVerify {
		options = append(options, klient.WithInsecureSkipVerify(true))
	}

	optionsRetry := []klient.OptionRetryFn{}
	if !cfg.Retry.Enabled {
		optionsRetry = append(optionsRetry, klient.OptionRetry.WithRetryDisable())
	}

	if len(cfg.Retry.DisabledStatusCodes) > 0 {
		optionsRetry = append(optionsRetry, klient.OptionRetry.WithRetryDisabledStatusCodes(cfg.Retry.DisabledStatusCodes...))
	}

	if len(cfg.Retry.EnabledStatusCodes) > 0 {
		optionsRetry = append(optionsRetry, klient.OptionRetry.WithRetryEnabledStatusCodes(cfg.Retry.EnabledStatusCodes...))
	}

	if len(optionsRetry) > 0 {
		options = append(options, klient.WithRetryOptions(optionsRetry...))
	}

	if cfg.Auth.Enabled {
		roundTripper, err := GlobalRegistry.AddService(cfg.Auth) //nolint:contextcheck // ctx from application level
		if err != nil {
			return nil, err
		}

		options = append(options, klient.WithRoundTripper(roundTripper))
	}

	client, err := klient.New(options...)
	if err != nil {
		return nil, err //nolint:wrapcheck // no need
	}

	return &Client{
		klient: client,
	}, nil
}

func (c *Client) Call(
	ctx context.Context,
	url, method string,
	headers map[string]interface{},
	payload []byte,
) (*ClientResponse, error) {
	req, err := c.newRequest(ctx, url, method, headers, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.klient.HTTP.Do(req)
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
func (c *Client) newRequest(ctx context.Context, url, method string, headers map[string]interface{}, payload []byte) (*http.Request, error) {
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
