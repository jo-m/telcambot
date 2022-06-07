.PHONY: run test

SHELL := /bin/bash

run:
	source .env && go run ./cmd/telcambot/ --log-pretty --log-level=debug --bot-debug=false

test:
	go test -v -count=1 ./...

telcambot-pi-zero:
	env \
		GOOS=linux \
		GOARCH=arm \
		GOARM=5 \
		CGO_ENABLED=0 \
	go build -o telcambot-pi-zero ./cmd/telcambot/
