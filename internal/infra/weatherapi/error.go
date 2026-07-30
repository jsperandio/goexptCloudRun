package weatherapi

import "errors"

var (
	ErrEmptyBaseURL       = errors.New("weatherapi base url is empty")
	ErrEmptyAPIKey        = errors.New("weatherapi api key is empty")
	ErrUnexpectedStatus   = errors.New("weatherapi returned an unexpected status")
	ErrMissingTemperature = errors.New("weatherapi response has no current temperature")
	ErrLocationMismatch   = errors.New("weatherapi resolved a different location")
)
