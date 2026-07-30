package viacep

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsperandio/goexptCloudRun/internal/entity"
)

const fullAddressPayload = `{
  "cep": "01001-000",
  "logradouro": "Praça da Sé",
  "complemento": "lado ímpar",
  "unidade": "",
  "bairro": "Sé",
  "localidade": "São Paulo",
  "uf": "SP",
  "estado": "São Paulo",
  "regiao": "Sudeste",
  "ibge": "3550308",
  "gia": "1004",
  "ddd": "11",
  "siafi": "7107"
}`

func Test_Response_ToLocation(t *testing.T) {
	t.Run("when the payload is the full upstream response, should decode every mirrored field", func(t *testing.T) {
		var payload Response
		require.NoError(t, json.Unmarshal([]byte(fullAddressPayload), &payload))

		assert.Equal(t, "01001-000", payload.Cep)
		assert.Equal(t, "Praça da Sé", payload.Logradouro)
		assert.Equal(t, "lado ímpar", payload.Complemento)
		assert.Equal(t, "Sé", payload.Bairro)
		assert.Equal(t, "São Paulo", payload.Localidade)
		assert.Equal(t, "SP", payload.Uf)
		assert.Equal(t, "São Paulo", payload.Estado)
		assert.Equal(t, "Sudeste", payload.Regiao)
		assert.Equal(t, "3550308", payload.Ibge)
		assert.Equal(t, "1004", payload.Gia)
		assert.Equal(t, "11", payload.Ddd)
		assert.Equal(t, "7107", payload.Siafi)
	})

	t.Run("when the payload has a locality, should map city, uf and state name", func(t *testing.T) {
		var payload Response
		require.NoError(t, json.Unmarshal([]byte(fullAddressPayload), &payload))

		location, err := payload.ToLocation()
		require.NoError(t, err)

		assert.Equal(t, "São Paulo", location.City)
		assert.Equal(t, "SP", location.UF)
		assert.Equal(t, "São Paulo", location.StateName)
	})

	t.Run("when the locality is empty, should return ErrZipcodeNotFound", func(t *testing.T) {
		var payload Response
		require.NoError(t, json.Unmarshal([]byte(`{"erro": "true"}`), &payload))

		_, err := payload.ToLocation()

		assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
	})
}
