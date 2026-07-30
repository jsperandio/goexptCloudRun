package viacep

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
	"github.com/jsperandio/goexptCloudRun/internal/infra/httpclient"
)

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

func (cl *Client) FindByZipcode(ctx context.Context, zipcode string) (entity.Location, error) {
	rq := NewRequest(zipcode)

	var rst Response
	resp, err := cl.client.R().
		SetContext(ctx).
		SetPathParam("zipcode", rq.Zipcode).
		ForceContentType(httpclient.ContentTypeJSON).
		SetResult(&rst).
		Get("/{zipcode}/json/")
	if err != nil {
		return entity.Location{}, fmt.Errorf("viacep request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return entity.Location{}, statusError(resp)
	}

	return rst.ToLocation()
}

func statusError(resp *resty.Response) error {
	return fmt.Errorf("%w: status %d: %s", ErrUnexpectedStatus, resp.StatusCode(), resp.String())
}
