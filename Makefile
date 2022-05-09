BINARY=main

lint:
	golangci-lint run ./...

swag:
	swag init -g cmd/api/main.go --parseDependency

build:
	go build -tags -o ${BINARY_HTTP} cmd/api/main.go