package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Temperature_Celsius(t *testing.T) {
	t.Run("when built from celsius, should return the same value", func(t *testing.T) {
		assert.Equal(t, 28.5, NewTemperatureFromCelsius(28.5).Celsius())
	})

	t.Run("when the value is negative, should keep the sign", func(t *testing.T) {
		assert.Equal(t, -40.0, NewTemperatureFromCelsius(-40).Celsius())
	})

	t.Run("when the value has more than two decimals, should round to two", func(t *testing.T) {
		assert.Equal(t, 21.35, NewTemperatureFromCelsius(21.3456).Celsius())
	})
}

func Test_Temperature_Fahrenheit(t *testing.T) {
	t.Run("when the temperature is 28.5, should return 83.3", func(t *testing.T) {
		assert.Equal(t, 83.3, NewTemperatureFromCelsius(28.5).Fahrenheit())
	})

	t.Run("when the temperature is zero, should return 32", func(t *testing.T) {
		assert.Equal(t, 32.0, NewTemperatureFromCelsius(0).Fahrenheit())
	})

	t.Run("when the temperature is 100, should return 212", func(t *testing.T) {
		assert.Equal(t, 212.0, NewTemperatureFromCelsius(100).Fahrenheit())
	})

	t.Run("when the temperature is -40, should return -40", func(t *testing.T) {
		assert.Equal(t, -40.0, NewTemperatureFromCelsius(-40).Fahrenheit())
	})
}

func Test_Temperature_Kelvin(t *testing.T) {
	t.Run("when the temperature is 28.5, should return 301.5", func(t *testing.T) {
		assert.Equal(t, 301.5, NewTemperatureFromCelsius(28.5).Kelvin())
	})

	t.Run("when the temperature is zero, should return 273", func(t *testing.T) {
		assert.Equal(t, 273.0, NewTemperatureFromCelsius(0).Kelvin())
	})

	t.Run("when the temperature is 100, should return 373", func(t *testing.T) {
		assert.Equal(t, 373.0, NewTemperatureFromCelsius(100).Kelvin())
	})

	t.Run("when the temperature is -40, should return 233", func(t *testing.T) {
		assert.Equal(t, 233.0, NewTemperatureFromCelsius(-40).Kelvin())
	})
}
