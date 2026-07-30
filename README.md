# goexptCloudRun

Serviço em Go que recebe um CEP brasileiro, identifica a cidade e devolve a temperatura atual em
Celsius, Fahrenheit e Kelvin.

## URL em produção

    https://goexpt-cloud-run-329954409839.us-central1.run.app

    curl -s https://goexpt-cloud-run-329954409839.us-central1.run.app/weather/01001000

## Endpoint desejado

`GET /weather/{cep}`, com `cep` sempre com 8 dígitos, sem hífen e sem espaço.

### Ex:
 Sucesso, 200:

    curl -s localhost:8080/weather/01001000
    {"temp_C":23.1,"temp_F":73.58,"temp_K":296.1}

CEP com formato inválido, 422:

    curl -s localhost:8080/weather/01001-000
    {"message":"invalid zipcode"}

CEP bem formado mas inexistente na base, 404:

    curl -s localhost:8080/weather/00000000
    {"message":"can not find zipcode"}

## Considerações sobre Kelvin

O enunciado dá a fórmula `K = C + 273`, mas o exemplo de resposta tem um valor diferente do que a fórmula produziria.
(`28.5 °C` para `301.65 K`) usa na prática `273.15`. 

Escolhido levar o exemplo erro e fixar na formula apresentada, `C + 273`, então `28.5 °C` sai como `301.5`, não `301.65`.

## Executando com Docker

    docker build -t clima-cep .
    docker run --rm -p 8080:8080 --env-file .env clima-cep

Ou usando o makefile:

    make docker-build
    make docker-run

O `.env` precisa ter no mínimo `WEATHER_API_KEY`. Sem ela o container encerra ao subir.

## Rodando os testes

    make test          
    make test-docker   # dentro de um container sem precisar de Go instalado

## Variáveis de ambiente

| Variável                 | Default                          | Obrigatória |
| ------------------------ | --------------------------------- | ----------- |
| `PORT`                   | `8080`                             | não         |
| `HTTP_GRACEFUL_TIMEOUT`  | `10s`                              | não         |
| `WEATHER_API_KEY`        |                                    | sim         |
| `WEATHER_API_BASE_URL`   | `https://api.weatherapi.com/v1`   | não         |
| `WEATHER_API_TIMEOUT`    | `3s`                               | não         |
| `WEATHER_API_RETRY_COUNT`| `1`                                | não         |
| `WEATHER_API_RETRY_WAIT` | `200ms`                            | não         |
| `VIACEP_BASE_URL`        | `https://viacep.com.br/ws`        | não         |
| `VIACEP_TIMEOUT`         | `3s`                               | não         |
| `VIACEP_RETRY_COUNT`     | `1`                                | não         |
| `VIACEP_RETRY_WAIT`      | `200ms`                            | não         |

## Documentação da API

A aplicação usa `swag` para gerar a documentação da API em OpenAPI 2.0 (Swagger). A documentação é
gerada em `docs/swagger.json` e `docs/swagger.yaml`, e a Swagger UI é servida na raiz (`/`, tanto em
http://localhost:8080/ local quanto na URL do Cloud Run acima): ela redireciona sozinha para a Swagger
UI em `/docs/index.html`, com os dois endpoints

    GET /weather/{cep}
    GET /health

