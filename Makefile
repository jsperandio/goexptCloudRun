.PHONY: run build test mocks swagger

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
