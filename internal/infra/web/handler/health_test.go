package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HealthHandler_Handle(t *testing.T) {
	t.Run("when called, should return ok", func(t *testing.T) {
		c, rec := newTestContext(t, "")

		require.NoError(t, NewHealthHandler().Handle(c))

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())
	})
}
