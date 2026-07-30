package weatherapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

const fullCurrentPayload = `{
  "location": {
    "name": "Sao Paulo",
    "region": "Sao Paulo",
    "country": "Brazil",
    "lat": -23.5333,
    "lon": -46.6167,
    "tz_id": "America/Sao_Paulo",
    "localtime_epoch": 1785387765,
    "localtime": "2026-07-30 02:02"
  },
  "current": {
    "last_updated_epoch": 1785387600,
    "last_updated": "2026-07-30 02:00",
    "temp_c": 17.5,
    "temp_f": 63.6,
    "is_day": 0,
    "condition": {
      "text": "Clear",
      "icon": "//cdn.weatherapi.com/weather/64x64/night/113.png",
      "code": 1000
    },
    "wind_mph": 2.2,
    "wind_kph": 3.6,
    "wind_degree": 163,
    "wind_dir": "SSE",
    "pressure_mb": 1019.0,
    "pressure_in": 30.09,
    "precip_mm": 0.0,
    "precip_in": 0.0,
    "humidity": 69,
    "cloud": 8,
    "feelslike_c": 12.9,
    "feelslike_f": 55.2,
    "windchill_c": 17.5,
    "windchill_f": 63.6,
    "heatindex_c": 17.5,
    "heatindex_f": 63.6,
    "dewpoint_c": 11.9,
    "dewpoint_f": 53.4,
    "vis_km": 10.0,
    "vis_miles": 6.0,
    "uv": 0.0,
    "gust_mph": 4.1,
    "gust_kph": 6.6,
    "will_it_rain": 0,
    "chance_of_rain": 6,
    "will_it_snow": 0,
    "chance_of_snow": 0,
    "wetbulb_c": 14.2,
    "wetbulb_f": 57.5,
    "short_rad": 0,
    "diff_rad": 0,
    "dni": 0,
    "gti": 0
  }
}`

func Test_Response_ToTemperature(t *testing.T) {
	t.Run("when the payload is the full upstream response, should decode every mirrored field", func(t *testing.T) {
		var payload Response
		require.NoError(t, json.Unmarshal([]byte(fullCurrentPayload), &payload))

		assert.Equal(t, "Sao Paulo", payload.Location.Name)
		assert.Equal(t, "Sao Paulo", payload.Location.Region)
		assert.Equal(t, "Brazil", payload.Location.Country)
		assert.Equal(t, -23.5333, payload.Location.Lat)
		assert.Equal(t, "America/Sao_Paulo", payload.Location.TzID)
		assert.Equal(t, int64(1785387765), payload.Location.LocaltimeEpoch)

		assert.Equal(t, "Clear", payload.Current.Condition.Text)
		assert.Equal(t, 1000, payload.Current.Condition.Code)
		assert.Equal(t, 69, payload.Current.Humidity)
		assert.Equal(t, "SSE", payload.Current.WindDir)
		assert.Equal(t, 12.9, payload.Current.FeelslikeC)
		assert.Equal(t, 6, payload.Current.ChanceOfRain)
		assert.Equal(t, 14.2, payload.Current.WetbulbC)
	})

	t.Run("when the payload has temp_c, should build the temperature from it", func(t *testing.T) {
		var payload Response
		require.NoError(t, json.Unmarshal([]byte(fullCurrentPayload), &payload))

		temperature, err := payload.ToTemperature()
		require.NoError(t, err)

		assert.Equal(t, 17.5, temperature.Celsius())
	})

	t.Run("when temp_c is absent, should return ErrMissingTemperature", func(t *testing.T) {
		var payload Response
		require.NoError(t, json.Unmarshal([]byte(`{"current":{"humidity":69}}`), &payload))

		_, err := payload.ToTemperature()

		assert.ErrorIs(t, err, ErrMissingTemperature)
	})

	t.Run("when temp_c is zero, should not confuse it with an absent value", func(t *testing.T) {
		var payload Response
		require.NoError(t, json.Unmarshal([]byte(`{"current":{"temp_c":0}}`), &payload))

		temperature, err := payload.ToTemperature()
		require.NoError(t, err)

		assert.Equal(t, 0.0, temperature.Celsius())
	})
}

func Test_ErrorResponse_ToError(t *testing.T) {
	t.Run("when the code is the location not found one, should return ErrZipcodeNotFound", func(t *testing.T) {
		var payload ErrorResponse
		body := `{"error":{"code":1006,"message":"No matching location found."}}`
		require.NoError(t, json.Unmarshal([]byte(body), &payload))

		assert.ErrorIs(t, payload.ToError(http.StatusBadRequest), entity.ErrZipcodeNotFound)
	})

	t.Run("when the code is any other, should report the status and the code", func(t *testing.T) {
		var payload ErrorResponse
		body := `{"error":{"code":2006,"message":"API key is invalid."}}`
		require.NoError(t, json.Unmarshal([]byte(body), &payload))

		err := payload.ToError(http.StatusUnauthorized)

		assert.ErrorIs(t, err, ErrUnexpectedStatus)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
		assert.Contains(t, err.Error(), "status 401")
		assert.Contains(t, err.Error(), "code 2006")
	})
}
