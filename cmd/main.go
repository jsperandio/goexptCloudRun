package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jsperandio/goexptCloudRun/configs"
	"github.com/jsperandio/goexptCloudRun/internal/infra/validator"
	"github.com/jsperandio/goexptCloudRun/internal/infra/viacep"
	"github.com/jsperandio/goexptCloudRun/internal/infra/weatherapi"
	"github.com/jsperandio/goexptCloudRun/internal/infra/web/handler"
	"github.com/jsperandio/goexptCloudRun/internal/infra/web/webserver"
	"github.com/jsperandio/goexptCloudRun/internal/usecase"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)))

	if err := run(); err != nil {
		slog.Error("application terminated", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped gracefully")
}

func run() error {
	configs.LoadEnv(".env")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	locationFinder, err := viacep.NewClient(nil)
	if err != nil {
		return err
	}

	weatherProvider, err := weatherapi.NewClient(nil)
	if err != nil {
		return err
	}

	getWeather := usecase.NewGetWeatherByZipcodeUseCase(locationFinder, weatherProvider, validator.NewValidator())

	ws, err := webserver.NewWebServer(nil)
	if err != nil {
		return err
	}

	weatherHandler := handler.NewWeatherHandler(getWeather)
	healthHandler := handler.NewHealthHandler()

	ws.RegisterRoute(http.MethodGet, handler.RouteWeather, weatherHandler.Handle)
	ws.RegisterRoute(http.MethodGet, handler.RouteHealth, healthHandler.Handle)

	slog.Info("starting application")

	return ws.Start(ctx)
}
