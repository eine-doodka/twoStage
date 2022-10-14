.PHONY: build
build:
	go build -v ./main.go

.PHONY: test
test:
	go test -v -race -timeout 10s ./...

.DEFAULT_GOAL := build