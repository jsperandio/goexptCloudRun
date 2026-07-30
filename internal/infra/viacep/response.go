package viacep

import "github.com/jsperandio/goexptCloudRun/internal/entity"

type Response struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func (rs Response) ToLocation() (entity.Location, error) {
	if rs.Localidade == "" {
		return entity.Location{}, entity.ErrZipcodeNotFound
	}

	return entity.Location{
		City:      rs.Localidade,
		UF:        rs.Uf,
		StateName: rs.Estado,
	}, nil
}
