package viacep

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type ClientOptions struct {
	BaseURL       string        `env:"VIACEP_BASE_URL" envDefault:"https://viacep.com.br/ws"`
	Timeout       time.Duration `env:"VIACEP_TIMEOUT" envDefault:"3s"`
	RetryCount    int           `env:"VIACEP_RETRY_COUNT" envDefault:"1"`
	RetryWaitTime time.Duration `env:"VIACEP_RETRY_WAIT" envDefault:"200ms"`
}

func NewDefaultClientOptions() (*ClientOptions, error) {
	opt, err := env.ParseAs[ClientOptions]()
	if err != nil {
		return nil, err
	}

	return &opt, nil
}
