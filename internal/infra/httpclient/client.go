package httpclient

import (
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	ContentTypeJSON = "application/json"
)

type Config struct {
	BaseURL       string
	Timeout       time.Duration
	RetryCount    int
	RetryWaitTime time.Duration
}

func New(cfg Config) *resty.Client {
	return resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout).
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(cfg.RetryWaitTime).
		AddRetryCondition(ShouldRetry).
		SetLogger(NoopLogger{})
}
