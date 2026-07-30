package handler

import (
	"errors"
	"net/http"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
	"github.com/jsperandio/goexptCloudRun/internal/usecase"
)

const internalErrorMessage = "internal server error"

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewWeatherResponseFromOutput(out usecase.GetWeatherByZipcodeOutput) WeatherResponse {
	return WeatherResponse{
		TempC: out.TempC,
		TempF: out.TempF,
		TempK: out.TempK,
	}
}

func NewErrorResponseFromError(err error) (int, ErrorResponse) {
	switch {
	case errors.Is(err, entity.ErrInvalidZipcode):
		return http.StatusUnprocessableEntity, ErrorResponse{
			Message: entity.ErrInvalidZipcode.Error(),
		}
	case errors.Is(err, entity.ErrZipcodeNotFound):
		return http.StatusNotFound, ErrorResponse{
			Message: entity.ErrZipcodeNotFound.Error(),
		}
	default:
		return http.StatusInternalServerError, ErrorResponse{
			Message: internalErrorMessage,
		}
	}
}
