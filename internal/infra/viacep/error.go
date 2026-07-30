package viacep

import "errors"

var (
	ErrEmptyBaseURL     = errors.New("viacep base url is empty")
	ErrUnexpectedStatus = errors.New("viacep returned an unexpected status")
)
