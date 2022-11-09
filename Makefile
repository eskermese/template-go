.PHONY:
.SILENT:
.DEFAULT_GOAL := all

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/app/main.go

run: build
	./.bin/app

swag:
	swag init -g cmd/app/main.go

lint:
	golangci-lint run

migrate:
	migrate -database postgres://postgres:qwerty123@localhost:5432/postgres?sslmode=disable -path migrations up

drop-tables:
	migrate -database postgres://postgres:qwerty123@localhost:5432/postgres?sslmode=disable -path migrations down
