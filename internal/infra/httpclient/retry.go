package httpclient

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

func ShouldRetry(resp *resty.Response, err error) bool {
	if err != nil {
		return resp == nil || resp.RawResponse == nil
	}

	return resp.StatusCode() >= http.StatusInternalServerError
}
