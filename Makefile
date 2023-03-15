.PHONY: build
build:
	GOOS=linux go build -o bin/loadBalancer cmd/loadBalancer/main.go