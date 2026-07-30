# goexptCloudRun

## Documentação da API

Com o servidor no ar, abra a URL base do serviço no navegador (`http://localhost:8080/` local, ou a
URL do Cloud Run) — ela redireciona sozinha para a Swagger UI em `/docs/index.html`, com os dois
endpoints (`GET /weather/{cep}` e `GET /health`) documentados.

As anotações `@Summary`/`@Param`/`@Success`/`@Router` em `cmd/main.go` e nos handlers de
`internal/infra/web/handler/` são uma exceção deliberada à convenção de "sem comentários no
código": elas não explicam o código, são o metadado estrutural que a ferramenta `swag` lê para
gerar a spec. Depois de mudar uma anotação, regenere os arquivos em `docs/` com:

    make swagger