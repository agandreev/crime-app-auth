BINARY=main

lint:
	golangci-lint run ./...

swag:
	swag init -g cmd/api/main.go --parseDependency

build:
	go build -tags -o ${BINARY_HTTP} cmd/api/main.go

vegeta:
	cd && echo "GET http://localhost" | vegeta attack -duration=30s -rate=10 -output=results-veg-httpbin-get.bin && cat results-veg-httpbin-get.bin | vegeta plot --title="HTTP Bin GET 10 rps for 30 seconds" > http-bin-get-10rps-30seconds.html