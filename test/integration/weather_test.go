package integration

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsperandio/goexptCloudRun/internal/infra/validator"
	"github.com/jsperandio/goexptCloudRun/internal/infra/viacep"
	"github.com/jsperandio/goexptCloudRun/internal/infra/weatherapi"
	"github.com/jsperandio/goexptCloudRun/internal/infra/web/handler"
	"github.com/jsperandio/goexptCloudRun/internal/usecase"
)

const (
	existingZipcode = "01001000"
	missingZipcode  = "99999999"
	testAPIKey      = "integration-key"
)

func newViaCEPServer(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if strings.HasPrefix(r.URL.Path, "/"+existingZipcode+"/") {
			_, _ = w.Write([]byte(`{"localidade":"São Paulo","uf":"SP","estado":"São Paulo"}`))

			return
		}

		_, _ = w.Write([]byte(`{"erro": "true"}`))
	}))

	t.Cleanup(server.Close)

	return server
}

func newWeatherAPIServer(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Query().Get("q") != "Sao Paulo,Sao Paulo,Brazil" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":{"code":1006,"message":"No matching location found."}}`))

			return
		}

		_, _ = w.Write([]byte(`{
			"location": {"name": "Sao Paulo", "region": "Sao Paulo", "country": "Brazil"},
			"current": {"temp_c": 28.5}
		}`))
	}))

	t.Cleanup(server.Close)

	return server
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	viaCEPServer := newViaCEPServer(t)
	weatherAPIServer := newWeatherAPIServer(t)

	finder, err := viacep.NewClient(&viacep.ClientOptions{
		BaseURL:       viaCEPServer.URL,
		Timeout:       5 * time.Second,
		RetryCount:    1,
		RetryWaitTime: time.Millisecond,
	})
	require.NoError(t, err)

	provider, err := weatherapi.NewClient(&weatherapi.ClientOptions{
		BaseURL:       weatherAPIServer.URL,
		APIKey:        testAPIKey,
		Timeout:       5 * time.Second,
		RetryCount:    1,
		RetryWaitTime: time.Millisecond,
	})
	require.NoError(t, err)

	weatherHandler := handler.NewWeatherHandler(
		usecase.NewGetWeatherByZipcodeUseCase(finder, provider, validator.NewValidator()),
	)
	healthHandler := handler.NewHealthHandler()

	e := echo.New()
	e.Logger = slog.New(slog.DiscardHandler)
	e.GET(handler.RouteWeather, weatherHandler.Handle)
	e.GET(handler.RouteHealth, healthHandler.Handle)

	server := httptest.NewServer(e)
	t.Cleanup(server.Close)

	return server
}

func get(t *testing.T, url string) (int, string) {
	t.Helper()

	resp, err := http.Get(url)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(body)
}

func Test_Weather(t *testing.T) {
	t.Run("when the zipcode exists, should return the three scales", func(t *testing.T) {
		server := newTestServer(t)

		status, body := get(t, server.URL+"/weather/"+existingZipcode)

		assert.Equal(t, http.StatusOK, status)
		assert.JSONEq(t, `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`, body)
	})

	t.Run("when the zipcode has 7 digits, should return 422", func(t *testing.T) {
		server := newTestServer(t)

		status, body := get(t, server.URL+"/weather/0100100")

		assert.Equal(t, http.StatusUnprocessableEntity, status)
		assert.JSONEq(t, `{"message":"invalid zipcode"}`, body)
	})

	t.Run("when the zipcode has a hyphen, should return 422", func(t *testing.T) {
		server := newTestServer(t)

		status, body := get(t, server.URL+"/weather/01001-000")

		assert.Equal(t, http.StatusUnprocessableEntity, status)
		assert.JSONEq(t, `{"message":"invalid zipcode"}`, body)
	})

	t.Run("when the zipcode does not exist, should return 404", func(t *testing.T) {
		server := newTestServer(t)

		status, body := get(t, server.URL+"/weather/"+missingZipcode)

		assert.Equal(t, http.StatusNotFound, status)
		assert.JSONEq(t, `{"message":"can not find zipcode"}`, body)
	})

	t.Run("when the health route is called, should return ok", func(t *testing.T) {
		server := newTestServer(t)

		status, body := get(t, server.URL+"/health")

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "ok", body)
	})

	t.Run("when the zipcode is missing from the path, should return the router 404", func(t *testing.T) {
		server := newTestServer(t)

		status, _ := get(t, server.URL+"/weather/")

		assert.Equal(t, http.StatusNotFound, status)
	})
}
