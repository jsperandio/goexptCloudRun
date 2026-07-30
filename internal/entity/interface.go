package entity

import "context"

type LocationFinder interface {
	FindByZipcode(ctx context.Context, zipcode string) (Location, error)
}

type WeatherProvider interface {
	CurrentByLocation(ctx context.Context, lc Location) (Temperature, error)
}

type Validator interface {
	Validate(value any) error
}
