package viacep

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

const addressPayload = `{
  "cep": "01001-000",
  "logradouro": "Praça da Sé",
  "bairro": "Sé",
  "localidade": "São Paulo",
  "uf": "SP",
  "estado": "São Paulo"
}`

const notFoundHTML = `<!DOCTYPE HTML>
<html lang="pt-br">
<head><title>ViaCEP 400</title></head>
<body><h1>400</h1><h2>Bad Request</h2></body>
</html>`

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()

	opt := &ClientOptions{
		BaseURL:       baseURL,
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
			Timeout: time.Second,
		}

		_, err := NewClient(opt)

		assert.ErrorIs(t, err, ErrEmptyBaseURL)
	})

	t.Run("when the options are nil, should read them from the environment", func(t *testing.T) {
		t.Setenv("VIACEP_BASE_URL", "http://localhost:1/ws")
		t.Setenv("VIACEP_TIMEOUT", "1s")

		cl, err := NewClient(nil)
		require.NoError(t, err)

		assert.Equal(t, "http://localhost:1/ws", cl.options.BaseURL)
		assert.Equal(t, time.Second, cl.options.Timeout)
	})

	t.Run("when the retry options are not set, should default to a single retry", func(t *testing.T) {
		t.Setenv("VIACEP_BASE_URL", "http://localhost:1/ws")

		cl, err := NewClient(nil)
		require.NoError(t, err)

		assert.Equal(t, 1, cl.options.RetryCount)
		assert.Equal(t, 200*time.Millisecond, cl.options.RetryWaitTime)
		assert.Equal(t, 3*time.Second, cl.options.Timeout)
	})
}

func Test_Client_FindByZipcode(t *testing.T) {
	t.Run("when the zipcode exists, should return the location", func(t *testing.T) {
		var requested string

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requested = r.URL.Path
			_, _ = w.Write([]byte(addressPayload))
		}))
		t.Cleanup(server.Close)

		location, err := newTestClient(t, server.URL).FindByZipcode(context.Background(), "01001000")
		require.NoError(t, err)

		assert.Equal(t, "São Paulo", location.City)
		assert.Equal(t, "SP", location.UF)
		assert.Equal(t, "São Paulo", location.StateName)
		assert.Equal(t, "/01001000/json/", requested)
	})

	t.Run("when the erro field is the string true, should return ErrZipcodeNotFound", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{"erro": "true"}`)

		_, err := newTestClient(t, server.URL).FindByZipcode(context.Background(), "99999999")

		assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the erro field is the boolean true, should return ErrZipcodeNotFound", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{"erro": true}`)

		_, err := newTestClient(t, server.URL).FindByZipcode(context.Background(), "99999999")

		assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the status is not ok, should not report the zipcode as not found", func(t *testing.T) {
		server := newServer(t, http.StatusBadRequest, notFoundHTML)

		_, err := newTestClient(t, server.URL).FindByZipcode(context.Background(), "123")

		assert.ErrorIs(t, err, ErrUnexpectedStatus)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
		assert.Contains(t, err.Error(), "status 400")
	})

	t.Run("when the payload is not valid json, should return an error", func(t *testing.T) {
		server := newServer(t, http.StatusOK, `{"localidade":`)

		_, err := newTestClient(t, server.URL).FindByZipcode(context.Background(), "01001000")

		assert.Error(t, err)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the first attempt fails with a server error, should retry once and return the location", func(t *testing.T) {
		server, calls := newCountingServer(t, func(w http.ResponseWriter, attempt int64) {
			if attempt == 1 {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, _ = w.Write([]byte(addressPayload))
		})

		location, err := newTestClient(t, server.URL).FindByZipcode(context.Background(), "01001000")
		require.NoError(t, err)

		assert.Equal(t, "São Paulo", location.City)
		assert.Equal(t, int64(2), calls.Load())
	})

	t.Run("when the status is a client error, should not retry", func(t *testing.T) {
		server, calls := newCountingServer(t, func(w http.ResponseWriter, _ int64) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(notFoundHTML))
		})

		_, err := newTestClient(t, server.URL).FindByZipcode(context.Background(), "123")

		assert.ErrorIs(t, err, ErrUnexpectedStatus)
		assert.Equal(t, int64(1), calls.Load())
	})

	t.Run("when the context is already cancelled, should return an error", func(t *testing.T) {
		server := newServer(t, http.StatusOK, addressPayload)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := newTestClient(t, server.URL).FindByZipcode(ctx, "01001000")

		assert.ErrorIs(t, err, context.Canceled)
	})
}
