package weatherapi

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

func bonitoIn(stateName string) entity.Location {
	return entity.Location{
		City:      "Bonito",
		StateName: stateName,
	}
}

func Test_Request_validateResolved(t *testing.T) {
	t.Run("when the resolved state is the requested one, should return nil", func(t *testing.T) {
		resolved := Location{
			Name:   "Bonito",
			Region: "Pernambuco",
		}

		err := NewRequest(bonitoIn("Pernambuco")).validateResolved(resolved)

		assert.NoError(t, err)
	})

	t.Run("when the resolved state differs, should return ErrLocationMismatch", func(t *testing.T) {
		resolved := Location{
			Name:   "Bonito",
			Region: "Mato Grosso do Sul",
		}

		err := NewRequest(bonitoIn("Pernambuco")).validateResolved(resolved)

		assert.ErrorIs(t, err, ErrLocationMismatch)
		assert.NotErrorIs(t, err, entity.ErrZipcodeNotFound)
		assert.Contains(t, err.Error(), "Pernambuco")
		assert.Contains(t, err.Error(), "Mato Grosso do Sul")
	})

	t.Run("when the upstream answers another country, should return ErrLocationMismatch", func(t *testing.T) {
		resolved := Location{
			Name:    "Brazil",
			Region:  "Indiana",
			Country: "United States of America",
		}

		location := entity.Location{
			City:      "Niterói",
			StateName: "Rio de Janeiro",
		}

		err := NewRequest(location).validateResolved(resolved)

		assert.ErrorIs(t, err, ErrLocationMismatch)
	})

	t.Run("when the resolved state differs only by accent, should accept it", func(t *testing.T) {
		resolved := Location{
			Name:   "Belem",
			Region: "Pará",
		}

		location := entity.Location{
			City:      "Belém",
			StateName: "Pará",
		}

		err := NewRequest(location).validateResolved(resolved)

		assert.NoError(t, err)
	})

	t.Run("when the resolved names differ only by case, should accept them", func(t *testing.T) {
		resolved := Location{
			Name:   "NITEROI",
			Region: "RIO DE JANEIRO",
		}

		location := entity.Location{
			City:      "Niterói",
			StateName: "Rio de Janeiro",
		}

		err := NewRequest(location).validateResolved(resolved)

		assert.NoError(t, err)
	})

	t.Run("when the city differs but the state matches, should return ErrLocationMismatch", func(t *testing.T) {
		resolved := Location{
			Name:   "Sao Paulo",
			Region: "Sao Paulo",
		}

		location := entity.Location{
			City:      "Embu",
			StateName: "São Paulo",
		}

		err := NewRequest(location).validateResolved(resolved)

		assert.ErrorIs(t, err, ErrLocationMismatch)
	})

	t.Run("when no state was requested, should still check the city", func(t *testing.T) {
		resolved := Location{
			Name:   "Sao Paulo",
			Region: "Sao Paulo",
		}

		err := NewRequest(bonitoIn("")).validateResolved(resolved)

		assert.ErrorIs(t, err, ErrLocationMismatch)
	})

	t.Run("when no state was requested and the city matches, should accept it", func(t *testing.T) {
		resolved := Location{
			Name:   "Bonito",
			Region: "Mato Grosso do Sul",
		}

		err := NewRequest(bonitoIn("")).validateResolved(resolved)

		assert.NoError(t, err)
	})
}
