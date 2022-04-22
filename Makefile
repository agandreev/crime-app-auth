lint:
	golangci-lint run ./...

swag:
	swag init -g cmd/api/main.go