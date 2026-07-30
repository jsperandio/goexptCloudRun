package entity

import "math"

const kelvinOffset = 273

type Temperature struct {
	celsius float64
}

func NewTemperatureFromCelsius(c float64) Temperature {
	return Temperature{
		celsius: c,
	}
}

func (tp Temperature) Celsius() float64 {
	return roundToTwoDecimals(tp.celsius)
}

func (tp Temperature) Fahrenheit() float64 {
	return roundToTwoDecimals(tp.celsius*1.8 + 32)
}

func (tp Temperature) Kelvin() float64 {
	return roundToTwoDecimals(tp.celsius + kelvinOffset)
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
