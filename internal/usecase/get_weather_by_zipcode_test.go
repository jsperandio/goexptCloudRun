package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
	"github.com/jsperandio/goexptCloudRun/internal/entity/mocks"
	"github.com/jsperandio/goexptCloudRun/internal/infra/validator"
)

func assertInvalidZipcode(t *testing.T, zipcode string) {
	t.Helper()

	lf := mocks.NewLocationFinder(t)
	wp := mocks.NewWeatherProvider(t)

	uc := NewGetWeatherByZipcodeUseCase(lf, wp, validator.NewValidator())

	in := GetWeatherByZipcodeInput{
		Zipcode: zipcode,
	}

	_, err := uc.Execute(context.Background(), in)

	assert.ErrorIs(t, err, entity.ErrInvalidZipcode)
}

func Test_GetWeatherByZipcodeUseCase_Execute(t *testing.T) {
	t.Run("when the zipcode has less than 8 digits, should return ErrInvalidZipcode", func(t *testing.T) {
		assertInvalidZipcode(t, "0100100")
	})

	t.Run("when the zipcode has more than 8 digits, should return ErrInvalidZipcode", func(t *testing.T) {
		assertInvalidZipcode(t, "010010000")
	})

	t.Run("when the zipcode is empty, should return ErrInvalidZipcode", func(t *testing.T) {
		assertInvalidZipcode(t, "")
	})

	t.Run("when the zipcode has a letter, should return ErrInvalidZipcode", func(t *testing.T) {
		assertInvalidZipcode(t, "0100100a")
	})

	t.Run("when the zipcode has a hyphen, should return ErrInvalidZipcode", func(t *testing.T) {
		assertInvalidZipcode(t, "01001-000")
	})

	t.Run("when the zipcode uses non ascii digits, should return ErrInvalidZipcode", func(t *testing.T) {
		assertInvalidZipcode(t, "٣٣٣٣٣٣٣٣")
	})

	t.Run("when the finder does not find the zipcode, should return ErrZipcodeNotFound", func(t *testing.T) {
		lf := mocks.NewLocationFinder(t)
		wp := mocks.NewWeatherProvider(t)

		lf.EXPECT().
			FindByZipcode(mock.Anything, "99999999").
			Return(entity.Location{}, entity.ErrZipcodeNotFound)

		uc := NewGetWeatherByZipcodeUseCase(lf, wp, validator.NewValidator())

		in := GetWeatherByZipcodeInput{
			Zipcode: "99999999",
		}

		_, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
	})

	t.Run("when the provider fails, should return the error", func(t *testing.T) {
		lf := mocks.NewLocationFinder(t)
		wp := mocks.NewWeatherProvider(t)

		location := entity.Location{
			City:      "São Paulo",
			UF:        "SP",
			StateName: "São Paulo",
		}

		providerErr := errors.New("upstream unavailable")

		lf.EXPECT().
			FindByZipcode(mock.Anything, "01001000").
			Return(location, nil)

		wp.EXPECT().
			CurrentByLocation(mock.Anything, location).
			Return(entity.Temperature{}, providerErr)

		uc := NewGetWeatherByZipcodeUseCase(lf, wp, validator.NewValidator())

		in := GetWeatherByZipcodeInput{
			Zipcode: "01001000",
		}

		_, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, providerErr)
	})

	t.Run("when the zipcode exists, should return the three scales", func(t *testing.T) {
		lf := mocks.NewLocationFinder(t)
		wp := mocks.NewWeatherProvider(t)

		location := entity.Location{
			City:      "São Paulo",
			UF:        "SP",
			StateName: "São Paulo",
		}

		lf.EXPECT().
			FindByZipcode(mock.Anything, "01001000").
			Return(location, nil)

		wp.EXPECT().
			CurrentByLocation(mock.Anything, location).
			Return(entity.NewTemperatureFromCelsius(28.5), nil)

		uc := NewGetWeatherByZipcodeUseCase(lf, wp, validator.NewValidator())

		in := GetWeatherByZipcodeInput{
			Zipcode: "01001000",
		}

		out, err := uc.Execute(context.Background(), in)
		require.NoError(t, err)

		assert.Equal(t, 28.5, out.TempC)
		assert.Equal(t, 83.3, out.TempF)
		assert.Equal(t, 301.5, out.TempK)
	})
}
