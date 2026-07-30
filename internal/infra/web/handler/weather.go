package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/jsperandio/goexptCloudRun/internal/usecase"
)

const RouteWeather string = "/weather/:cep"

type WeatherHandler struct {
	usecase *usecase.GetWeatherByZipcodeUseCase
}

func NewWeatherHandler(uc *usecase.GetWeatherByZipcodeUseCase) *WeatherHandler {
	return &WeatherHandler{
		usecase: uc,
	}
}

func (wh *WeatherHandler) Handle(c *echo.Context) error {
	in := usecase.GetWeatherByZipcodeInput{
		Zipcode: c.Param("cep"),
	}

	out, err := wh.usecase.Execute(c.Request().Context(), in)
	if err != nil {
		return wh.respondError(c, err)
	}

	return c.JSON(http.StatusOK, NewWeatherResponseFromOutput(out))
}

func (wh *WeatherHandler) respondError(c *echo.Context, err error) error {
	status, payload := NewErrorResponseFromError(err)

	if status == http.StatusInternalServerError {
		c.Logger().Error("weather request failed", "error", err)
	}

	return c.JSON(status, payload)
}
