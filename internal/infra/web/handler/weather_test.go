package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
	"github.com/jsperandio/goexptCloudRun/internal/entity/mocks"
	"github.com/jsperandio/goexptCloudRun/internal/infra/validator"
	"github.com/jsperandio/goexptCloudRun/internal/usecase"
)

func newTestContext(t *testing.T, zipcode string) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/weather/"+zipcode, nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetLogger(slog.New(slog.DiscardHandler))
	c.SetPathValues(echo.PathValues{
		{
			Name:  "cep",
			Value: zipcode,
		},
	})

	return c, rec
}

func newTestHandler(t *testing.T) (*WeatherHandler, *mocks.LocationFinder, *mocks.WeatherProvider) {
	t.Helper()

	lf := mocks.NewLocationFinder(t)
	wp := mocks.NewWeatherProvider(t)

	return NewWeatherHandler(usecase.NewGetWeatherByZipcodeUseCase(lf, wp, validator.NewValidator())), lf, wp
}

func saoPaulo() entity.Location {
	return entity.Location{
		City:      "São Paulo",
		UF:        "SP",
		StateName: "São Paulo",
	}
}

func Test_WeatherHandler_Weather(t *testing.T) {
	t.Run("when the zipcode exists, should return the three scales", func(t *testing.T) {
		wh, lf, wp := newTestHandler(t)

		lf.EXPECT().
			FindByZipcode(mock.Anything, "01001000").
			Return(saoPaulo(), nil)

		wp.EXPECT().
			CurrentByLocation(mock.Anything, saoPaulo()).
			Return(entity.NewTemperatureFromCelsius(28.5), nil)

		c, rec := newTestContext(t, "01001000")

		require.NoError(t, wh.Handle(c))

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}`, rec.Body.String())
	})

	t.Run("when the zipcode is invalid, should return 422 without calling the upstreams", func(t *testing.T) {
		wh, _, _ := newTestHandler(t)

		c, rec := newTestContext(t, "0100100")

		require.NoError(t, wh.Handle(c))

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.JSONEq(t, `{"message":"invalid zipcode"}`, rec.Body.String())
	})

	t.Run("when the zipcode does not exist, should return 404", func(t *testing.T) {
		wh, lf, _ := newTestHandler(t)

		lf.EXPECT().
			FindByZipcode(mock.Anything, "99999999").
			Return(entity.Location{}, entity.ErrZipcodeNotFound)

		c, rec := newTestContext(t, "99999999")

		require.NoError(t, wh.Handle(c))

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"message":"can not find zipcode"}`, rec.Body.String())
	})

	t.Run("when the provider fails, should return 500 without leaking the upstream detail", func(t *testing.T) {
		wh, lf, wp := newTestHandler(t)

		providerErr := errors.New("weatherapi returned an unexpected status: status 401: code 2006")

		lf.EXPECT().
			FindByZipcode(mock.Anything, "01001000").
			Return(saoPaulo(), nil)

		wp.EXPECT().
			CurrentByLocation(mock.Anything, saoPaulo()).
			Return(entity.Temperature{}, providerErr)

		c, rec := newTestContext(t, "01001000")

		require.NoError(t, wh.Handle(c))

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"message":"internal server error"}`, rec.Body.String())
		assert.NotContains(t, rec.Body.String(), "2006")
		assert.NotContains(t, rec.Body.String(), "weatherapi")
	})
}
