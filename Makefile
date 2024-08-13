.PHONY: build clean

GOPATH ?= ~/go

all: build

build:
	mkdir -p build
	go build -o build/server main.go

clean:
	rm -rf build

test:
	${GOPATH}/bin/golangci-lint run
	go test ./...
