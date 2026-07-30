package httpclient_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jsperandio/goexptCloudRun/internal/infra/httpclient"
)

func Test_New(t *testing.T) {
	t.Run("when built, should apply the base url and retry settings", func(t *testing.T) {
		cl := httpclient.New(httpclient.Config{
			BaseURL:       "http://localhost:1",
			Timeout:       2 * time.Second,
			RetryCount:    3,
			RetryWaitTime: 50 * time.Millisecond,
		})

		assert.Equal(t, "http://localhost:1", cl.BaseURL)
		assert.Equal(t, 3, cl.RetryCount)
		assert.Equal(t, 50*time.Millisecond, cl.RetryWaitTime)
	})
}
