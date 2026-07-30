package weatherapi

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

func Test_Request_Query(t *testing.T) {
	t.Run("when the city and the state have accents, should send them without accents", func(t *testing.T) {
		location := entity.Location{
			City:      "Niterói",
			UF:        "RJ",
			StateName: "Rio de Janeiro",
		}

		query := NewRequest(location).Query(testAPIKey)

		assert.Equal(t, "Niteroi,Rio de Janeiro,Brazil", query.Get("q"))
	})

	t.Run("when the state name itself has accents, should strip them too", func(t *testing.T) {
		location := entity.Location{
			City:      "Belém",
			UF:        "PA",
			StateName: "Pará",
		}

		query := NewRequest(location).Query(testAPIKey)

		assert.Equal(t, "Belem,Para,Brazil", query.Get("q"))
	})

	t.Run("when the city has a cedilla, should replace it as expected by the upstream", func(t *testing.T) {
		location := entity.Location{
			City:      "São Lourenço",
			UF:        "MG",
			StateName: "Minas Gerais",
		}

		query := NewRequest(location).Query(testAPIKey)

		assert.Equal(t, "Sao Lourenco,Minas Gerais,Brazil", query.Get("q"))
	})

	t.Run("when two cities share the name, should differ only by the state", func(t *testing.T) {
		bonitoMS := entity.Location{
			City:      "Bonito",
			UF:        "MS",
			StateName: "Mato Grosso do Sul",
		}

		bonitoPE := entity.Location{
			City:      "Bonito",
			UF:        "PE",
			StateName: "Pernambuco",
		}

		queryMS := NewRequest(bonitoMS).Query(testAPIKey)
		queryPE := NewRequest(bonitoPE).Query(testAPIKey)

		assert.Equal(t, "Bonito,Mato Grosso do Sul,Brazil", queryMS.Get("q"))
		assert.Equal(t, "Bonito,Pernambuco,Brazil", queryPE.Get("q"))
		assert.NotEqual(t, queryMS.Get("q"), queryPE.Get("q"))
	})

	t.Run("when the state name is missing, should fall back to city and country", func(t *testing.T) {
		location := entity.Location{
			City:      "Bonito",
			UF:        "MS",
			StateName: "",
		}

		query := NewRequest(location).Query(testAPIKey)

		assert.Equal(t, "Bonito,Brazil", query.Get("q"))
	})

	t.Run("when the query is built, should always ask for the current weather without air quality", func(t *testing.T) {
		query := NewRequest(saoPaulo()).Query(testAPIKey)

		assert.Equal(t, "no", query.Get("aqi"))
		assert.Equal(t, testAPIKey, query.Get("key"))
	})
}
