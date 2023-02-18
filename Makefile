.PHONY: run test build-raspi-zero

SHELL := /bin/bash

run:
	source .env && go run ./cmd/telcambot/ --log-pretty --log-level=debug --bot-debug=false

build-raspi-zero:
	env \
		GOOS=linux \
		GOARCH=arm \
		GOARM=5 \
		CGO_ENABLED=0 \
	go build -o telcambot-pi-zero ./cmd/telcambot/

lint:
	gofmt -l .; test -z "$$(gofmt -l .)"
	find . \( -name '*.c' -or -name '*.h' \) -exec clang-format-10 --style=file --dry-run --Werror {} +
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000 ./...
	go run github.com/securego/gosec/v2/cmd/gosec@latest ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

check: lint

format:
	gofmt -w -s .
	find . \( -name '*.c' -or -name '*.h' \) -exec clang-format-10 --style=file -i {} +
