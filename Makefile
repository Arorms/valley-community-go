.PHONY: css watch-css build run test fmt vet migrate-up migrate-down migrate-version

css:
	npm run css

watch-css:
	npm run watch:css

build:
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate

run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

migrate-up:
	go run ./cmd/migrate -dir up

migrate-down:
	go run ./cmd/migrate -dir down -steps 1

migrate-version:
	go run ./cmd/migrate -dir version
