SHELL := /bin/bash

BINARY := onepw-pdf-export

all: build

build:
	go build -o $(BINARY) ./

clean:
	rm -f $(BINARY)

test:
	go test ./...
