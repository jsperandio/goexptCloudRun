package usecase

import (
	"context"
	"fmt"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

type GetWeatherByZipcodeUseCase struct {
	finder    entity.LocationFinder
	provider  entity.Weather
	validator entity.Validator
}

func NewGetWeatherByZipcodeUseCase(lf entity.LocationFinder, wp entity.Weather, vl entity.Validator) *GetWeatherByZipcodeUseCase {
	return &GetWeatherByZipcodeUseCase{
		finder:    lf,
		provider:  wp,
		validator: vl,
	}
}

func (uc *GetWeatherByZipcodeUseCase) Execute(ctx context.Context, in GetWeatherByZipcodeInput) (GetWeatherByZipcodeOutput, error) {
	if err := uc.validator.Validate(in); err != nil {
		return GetWeatherByZipcodeOutput{}, fmt.Errorf("%w: %v", entity.ErrInvalidZipcode, err)
	}

	loc, err := uc.finder.FindByZipcode(ctx, in.Zipcode)
	if err != nil {
		return GetWeatherByZipcodeOutput{}, err
	}

	t, err := uc.provider.CurrentByLocation(ctx, loc)
	if err != nil {
		return GetWeatherByZipcodeOutput{}, err
	}

	return GetWeatherByZipcodeOutput{
		TempC: t.Celsius(),
		TempF: t.Fahrenheit(),
		TempK: t.Kelvin(),
	}, nil
}
