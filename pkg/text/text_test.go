package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsperandio/goexptCloudRun/pkg/text"
)

func Test_StripAccents(t *testing.T) {
	t.Run("when the value has no accents, should return it unchanged", func(t *testing.T) {
		assert.Equal(t, "Curitiba", text.StripAccents("Curitiba"))
	})

	t.Run("when the value has every portuguese diacritic, should remove all of them", func(t *testing.T) {
		assert.Equal(t, "aaaaeeiooouuc", text.StripAccents("áàâãéêíóôõúüç"))
	})

	t.Run("when the value is empty, should return empty", func(t *testing.T) {
		assert.Equal(t, "", text.StripAccents(""))
	})
}

func Test_EqualIgnoringAccents(t *testing.T) {
	t.Run("when the values differ only by accent, should be equal", func(t *testing.T) {
		assert.True(t, text.EqualIgnoringAccents("Belém", "Belem"))
	})

	t.Run("when the values differ only by case, should be equal", func(t *testing.T) {
		assert.True(t, text.EqualIgnoringAccents("NITEROI", "niteroi"))
	})

	t.Run("when the values differ beyond accent and case, should not be equal", func(t *testing.T) {
		assert.False(t, text.EqualIgnoringAccents("São Paulo", "Rio de Janeiro"))
	})
}
