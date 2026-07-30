package weatherapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

const testAPIKey = "s3cr3t-api-key"

const saoPauloLocation = `"location":{"name":"Sao Paulo","region":"Sao Paulo","country":"Brazil"}`

func saoPaulo() entity.Location {
	return entity.Location{
		City:      "São Paulo",
		UF:        "SP",
		StateName: "São Paulo",
	}
}

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()

	opt := &ClientOptions{
		BaseURL:       baseURL,
		APIKey:        testAPIKey,
		Timeout:       5 * time.Second,
		RetryCount:    1,
		RetryWaitTime: time.Millisecond,
	}

	cl, err := NewClient(opt)
	require.NoError(t, err)

	return cl
}

func newCountingServer(t *testing.T, handler func(w http.ResponseWriter, attempt int64)) (*httptest.Server, *atomic.Int64) {
	t.Helper()

	var calls atomic.Int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, calls.Add(1))
	}))

	t.Cleanup(server.Close)

	return server, &calls
}

func newServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))

	t.Cleanup(server.Close)

	return server
}

func Test_NewClient(t *testing.T) {
	t.Run("when the base url is empty, should return ErrEmptyBaseURL", func(t *testing.T) {
		opt := &ClientOptions{
			BaseURL: "",
			APIKey:  testAPIKey,
			Timeout: time.Second,
		}

		_, err := NewClient(opt)

		assert.ErrorIs(t, err, ErrEmptyBaseURL)
	})

	t.Run("when the api key is empty, should return ErrEmptyAPIKey", func(t *testing.T) {
		opt := &ClientOptions{
			BaseURL: "http://localhost:1/v1",
			APIKey:  "",
			Timeout: time.Second,
		}

		_, err := NewClient(opt)

		assert.ErrorIs(t, err, ErrEmptyAPIKey)
	})

	t.Run("when the options are nil, should read them from the environment", func(t *testing.T) {
		t.Setenv("WEATHER_API_BASE_URL", "http://localhost:1/v1")
		t.Setenv("WEATHER_API_KEY", testAPIKey)
		t.Setenv("WEATHER_API_TIMEOUT", "1s")

		cl, err := NewClient(nil)
		require.NoError(t, err)

		assert.Equal(t, "http://localhost:1/v1", cl.options.BaseURL)
		assert.Equal(t, testAPIKey, cl.options.APIKey)
		assert.Equal(t, time.Second, cl.options.Timeout)
	})

	t.Run("when the retry options are not set, should default to a single retry", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", testAPIKey)

		cl, err := NewClient(nil)
		require.NoError(t, err)

		assert.Equal(t, 1, cl.options.RetryCount)
		assert.Equal(t, 200*time.Millisecond, cl.options.RetryWaitTime)
		assert.Equal(t, 3*time.Second, cl.options.Timeout)
	})

	t.Run("when the api key is not set, should fail to build the options", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", testAPIKey)
		require.NoError(t, os.Unsetenv("WEATHER_API_KEY"))

		_, err := NewClient(nil)

		assert.Error(t, err)
	})

	t.Run("when the api key is set but empty, should fail to build the options", func(t *testing.T) {
		t.Setenv("WEATHER_API_KEY", "")

		_, err := NewClient(nil)

		assert.Error(t, err)
	})
}

func Test_Client_CurrentByLocation(t *testing.T) {
	t.Run("when the location has weather, should return the temperature", func(t *testing.T) {
		var requestedPath string
		var requestedQuery url.Values

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestedPath = r.URL.Path
			requestedQuery = r.URL.Query()

			_, _ = w.Write([]byte(`{` + saoPauloLocation + `,"current":{"temp_c":28.5}}`))
		}))
		t.Cleanup(server.Close)

		temperature, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())
		require.NoError(t, err)

		assert.Equal(t, 28.5, temperature.Celsius())
		assert.Equal(t, "/current.json", requestedPath)
		assert.Equal(t, "Sao Paulo,Sao Paulo,Brazil", requestedQuery.Get("q"))
		assert.Equal(t, testAPIKey, requestedQuery.Get("key"))
	})

	t.Run("when the temperature is zero, should return zero instead of an error", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{`+saoPauloLocation+`,"current":{"temp_c":0}}`)

		temperature, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())
		require.NoError(t, err)

		assert.Equal(t, 0.0, temperature.Celsius())
	})

	t.Run("when the location is not found, should return ErrZipcodeNotFound", func(t *testing.T) {
		body := `{"error":{"code":1006,"message":"No matching location found."}}`
		server := newServer(t, http.StatusBadRequest, body)

		_, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())

		assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the api key is invalid, should not report the zipcode as not found", func(t *testing.T) {
		body := `{"error":{"code":2006,"message":"API key is invalid."}}`
		server := newServer(t, http.StatusUnauthorized, body)

		_, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())

		assert.ErrorIs(t, err, ErrUnexpectedStatus)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the quota is exceeded, should not report the zipcode as not found", func(t *testing.T) {
		body := `{"error":{"code":2007,"message":"API key has exceeded calls per month quota."}}`
		server := newServer(t, http.StatusForbidden, body)

		_, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())

		assert.ErrorIs(t, err, ErrUnexpectedStatus)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when an error status has a malformed body, should return a generic error", func(t *testing.T) {
		server := newServer(t, http.StatusUnauthorized, `{`)

		_, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())

		assert.ErrorIs(t, err, ErrUnexpectedStatus)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the response has no current temperature, should return ErrMissingTemperature", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{`+saoPauloLocation+`}`)

		_, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())

		assert.ErrorIs(t, err, ErrMissingTemperature)
	})

	t.Run("when the payload is not valid json, should return an error", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{"current":`)

		_, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())

		assert.Error(t, err)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the first attempt fails with a server error, should retry once and return the temperature", func(t *testing.T) {
		server, calls := newCountingServer(t, func(w http.ResponseWriter, attempt int64) {
			if attempt == 1 {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, _ = w.Write([]byte(`{` + saoPauloLocation + `,"current":{"temp_c":28.5}}`))
		})

		temperature, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())
		require.NoError(t, err)

		assert.Equal(t, 28.5, temperature.Celsius())
		assert.Equal(t, int64(2), calls.Load())
	})

	t.Run("when the location is not found, should not retry", func(t *testing.T) {
		server, calls := newCountingServer(t, func(w http.ResponseWriter, _ int64) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":1006,"message":"No matching location found."}}`))
		})

		_, err := newTestClient(t, server.URL).CurrentByLocation(context.Background(), saoPaulo())

		assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
		assert.Equal(t, int64(1), calls.Load())
	})

	t.Run("when the request fails, should not leak the api key in the error", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{`+saoPauloLocation+`,"current":{"temp_c":28.5}}`)
		cl := newTestClient(t, server.URL)
		server.Close()

		_, err := cl.CurrentByLocation(context.Background(), saoPaulo())
		require.Error(t, err)

		assert.NotContains(t, err.Error(), testAPIKey)
	})

	t.Run("when the context is already cancelled, should return an error", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{`+saoPauloLocation+`,"current":{"temp_c":28.5}}`)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := newTestClient(t, server.URL).CurrentByLocation(ctx, saoPaulo())

		assert.ErrorIs(t, err, context.Canceled)
		assert.NotContains(t, err.Error(), testAPIKey)
	})
}
