.PHONY: run build test mocks swagger docker-build docker-run test-docker

run:
	go run ./cmd

build:
	go build ./...

test:
	go test ./... -race -count=1

mocks:
	go tool mockery

swagger:
	go tool swag init -g cmd/main.go -o docs

docker-build:
	docker build -t clima-cep .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env clima-cep

test-docker:
	docker run --rm -v $(PWD):/src -w /src -v go-mod-cache:/go/pkg/mod golang:1.26 go test ./... -race -count=1
