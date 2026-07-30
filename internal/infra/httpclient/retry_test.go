package httpclient_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/jsperandio/goexptCloudRun/internal/infra/httpclient"
)

func Test_ShouldRetry(t *testing.T) {
	t.Run("when there is a transport error and no response, should retry", func(t *testing.T) {
		assert.True(t, httpclient.ShouldRetry(nil, errors.New("boom")))
	})

	t.Run("when there is a transport error and the response has no raw response, should retry", func(t *testing.T) {
		resp := &resty.Response{}

		assert.True(t, httpclient.ShouldRetry(resp, errors.New("boom")))
	})

	t.Run("when the status is a server error, should retry", func(t *testing.T) {
		resp := &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusInternalServerError}}

		assert.True(t, httpclient.ShouldRetry(resp, nil))
	})

	t.Run("when the status is a client error, should not retry", func(t *testing.T) {
		resp := &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusBadRequest}}

		assert.False(t, httpclient.ShouldRetry(resp, nil))
	})

	t.Run("when the status is success, should not retry", func(t *testing.T) {
		resp := &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusOK}}

		assert.False(t, httpclient.ShouldRetry(resp, nil))
	})
}
