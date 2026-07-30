package float_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsperandio/goexptCloudRun/pkg/float"
)

func Test_RoundToTwoDecimals(t *testing.T) {
	t.Run("when the value already has two decimals, should return it unchanged", func(t *testing.T) {
		assert.Equal(t, 28.5, float.RoundToTwoDecimals(28.5))
	})

	t.Run("when the value has more than two decimals, should round to two", func(t *testing.T) {
		assert.Equal(t, 301.65, float.RoundToTwoDecimals(301.654999))
	})

	t.Run("when the value rounds up at the third decimal, should round up", func(t *testing.T) {
		assert.Equal(t, 83.3, float.RoundToTwoDecimals(83.295))
	})

	t.Run("when the value is negative, should round preserving the sign", func(t *testing.T) {
		assert.Equal(t, -12.35, float.RoundToTwoDecimals(-12.346))
	})

	t.Run("when the value is zero, should return zero", func(t *testing.T) {
		assert.Equal(t, 0.0, float.RoundToTwoDecimals(0))
	})
}
