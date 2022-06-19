package request

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type RetryCodesType string

const RetryCodesValue RetryCodesType = "retryCodeValue"

var ResponseErrLimit int64 = 1024

type Retry struct {
	DisabledStatusCodes []int
	EnabledStatusCodes  []int
}

// RetryPolicy provides a default callback for Client.CheckRetry, which
// will retry on connection errors and server errors.
func RetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	retryCodes, _ := ctx.Value(RetryCodesValue).(*Retry)

	if retryCodes != nil {
		for _, disabledStatusCode := range retryCodes.DisabledStatusCodes {
			if resp.StatusCode == disabledStatusCode {
				return false, nil
			}
		}

		for _, enabledStatusCode := range retryCodes.EnabledStatusCodes {
			if resp.StatusCode == enabledStatusCode {
				return true, fmt.Errorf("force retried HTTP status %s: [%s]", resp.Status, AddSomeResponse(resp))
			}
		}
	}

	v, err := retryablehttp.ErrorPropagatedRetryPolicy(ctx, resp, err)
	if v && err != nil {
		err = fmt.Errorf("%w: [%s]", err, AddSomeResponse(resp))
	}

	return v, err
}

// AddSomeResponse not close body, retry library draining it.
func AddSomeResponse(resp *http.Response) []byte {
	v, _ := io.ReadAll(io.LimitReader(resp.Body, ResponseErrLimit))

	return v
}
