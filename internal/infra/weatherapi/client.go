package weatherapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
	"github.com/jsperandio/goexptCloudRun/internal/infra/httpclient"
)

const currentPath = "/current.json"

type Client struct {
	options *ClientOptions
	client  *resty.Client
}

func NewClient(opt *ClientOptions) (*Client, error) {
	if opt == nil {
		parsed, err := NewDefaultClientOptions()
		if err != nil {
			return nil, err
		}

		opt = parsed
	}

	if opt.BaseURL == "" {
		return nil, ErrEmptyBaseURL
	}

	if opt.APIKey == "" {
		return nil, ErrEmptyAPIKey
	}

	return &Client{
		options: opt,
		client: httpclient.New(httpclient.Config{
			BaseURL:       opt.BaseURL,
			Timeout:       opt.Timeout,
			RetryCount:    opt.RetryCount,
			RetryWaitTime: opt.RetryWaitTime,
		}),
	}, nil
}

func (cl *Client) CurrentByLocation(ctx context.Context, lc entity.Location) (entity.Temperature, error) {
	rq := NewRequest(lc)

	var rst Response
	resp, err := cl.client.R().
		SetContext(ctx).
		SetQueryParamsFromValues(rq.Query(cl.options.APIKey)).
		ForceContentType(httpclient.ContentTypeJSON).
		SetResult(&rst).
		SetError(&ErrorResponse{}).
		Get(currentPath)
	if err != nil {
		return entity.Temperature{}, fmt.Errorf("weatherapi request failed: %w", withoutURL(err))
	}

	if resp.StatusCode() != http.StatusOK {
		return entity.Temperature{}, statusError(resp)
	}

	if err := rq.validateResolved(rst.Location); err != nil {
		return entity.Temperature{}, err
	}

	return rst.ToTemperature()
}

func statusError(resp *resty.Response) error {
	payload := resp.Error().(*ErrorResponse)
	if payload.Error.Code == 0 {
		return fmt.Errorf("%w: status %d", ErrUnexpectedStatus, resp.StatusCode())
	}

	return payload.ToError(resp.StatusCode())
}

func withoutURL(err error) error {
	if urlErr, ok := errors.AsType[*url.Error](err); ok {
		return urlErr.Err
	}

	return err
}
