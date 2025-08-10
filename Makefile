.DEFAULT_GOAL := build
.PHONY: fmt lint vet build clean

fmt:
	go fmt ./...

lint: fmt
	golint ./...

vet: fmt
	go vet ./...

build: vet
	go build -o cmd/main cmd/main.go 

clean:
	rm cmd/main

