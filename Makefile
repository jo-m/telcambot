.PHONY: run test

SHELL := /bin/bash

run:
	source .env && go run ./cmd/telcambot/ --log-pretty --log-level=debug --bot-debug=false

test:
	go test -v -count=1 ./...
