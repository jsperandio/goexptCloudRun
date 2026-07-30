package weatherapi

import (
	"fmt"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

const locationNotFoundCode = 1006

type Response struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Current struct {
	LastUpdatedEpoch int64     `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            *float64  `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMph          float64   `json:"wind_mph"`
	WindKph          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMb       float64   `json:"pressure_mb"`
	PressureIn       float64   `json:"pressure_in"`
	PrecipMm         float64   `json:"precip_mm"`
	PrecipIn         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	WindchillC       float64   `json:"windchill_c"`
	WindchillF       float64   `json:"windchill_f"`
	HeatindexC       float64   `json:"heatindex_c"`
	HeatindexF       float64   `json:"heatindex_f"`
	DewpointC        float64   `json:"dewpoint_c"`
	DewpointF        float64   `json:"dewpoint_f"`
	VisKm            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	UV               float64   `json:"uv"`
	GustMph          float64   `json:"gust_mph"`
	GustKph          float64   `json:"gust_kph"`
	WillItRain       int       `json:"will_it_rain"`
	ChanceOfRain     int       `json:"chance_of_rain"`
	WillItSnow       int       `json:"will_it_snow"`
	ChanceOfSnow     int       `json:"chance_of_snow"`
	WetbulbC         float64   `json:"wetbulb_c"`
	WetbulbF         float64   `json:"wetbulb_f"`
	ShortRad         float64   `json:"short_rad"`
	DiffRad          float64   `json:"diff_rad"`
	Dni              float64   `json:"dni"`
	Gti              float64   `json:"gti"`
}

func (rs Response) ToTemperature() (entity.Temperature, error) {
	if rs.Current.TempC == nil {
		return entity.Temperature{}, ErrMissingTemperature
	}

	return entity.NewTemperatureFromCelsius(*rs.Current.TempC), nil
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (rs ErrorResponse) ToError(status int) error {
	if rs.Error.Code == locationNotFoundCode {
		return entity.ErrZipcodeNotFound
	}

	return fmt.Errorf(
		"%w: status %d: code %d: %s",
		ErrUnexpectedStatus,
		status,
		rs.Error.Code,
		rs.Error.Message,
	)
}
