package weatherapi

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type ClientOptions struct {
	BaseURL       string        `env:"WEATHER_API_BASE_URL" envDefault:"https://api.weatherapi.com/v1"`
	APIKey        string        `env:"WEATHER_API_KEY,required,notEmpty"`
	Timeout       time.Duration `env:"WEATHER_API_TIMEOUT" envDefault:"3s"`
	RetryCount    int           `env:"WEATHER_API_RETRY_COUNT" envDefault:"1"`
	RetryWaitTime time.Duration `env:"WEATHER_API_RETRY_WAIT" envDefault:"200ms"`
}

func NewDefaultClientOptions() (*ClientOptions, error) {
	opt, err := env.ParseAs[ClientOptions]()
	if err != nil {
		return nil, err
	}

	return &opt, nil
}
