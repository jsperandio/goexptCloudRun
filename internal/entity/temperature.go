package entity

import "github.com/jsperandio/goexptCloudRun/pkg/float"

type Temperature struct {
	celsius float64
}

func NewTemperatureFromCelsius(c float64) Temperature {
	return Temperature{
		celsius: c,
	}
}

func (tp Temperature) Celsius() float64 {
	return float.RoundToTwoDecimals(tp.celsius)
}

// Celsius para Fahrenheit: F = C × 1.8 + 32
func (tp Temperature) Fahrenheit() float64 {
	return float.RoundToTwoDecimals(tp.celsius*1.8 + 32)
}

// Celsius para Kelvin: K = C + 273
func (tp Temperature) Kelvin() float64 {
	return float.RoundToTwoDecimals(tp.celsius + 273)
}
