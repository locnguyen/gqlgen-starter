include .env
export $(shell sed 's/=.*//' .env)

versionFile := version.txt
VERSION := $(shell cat ${versionFile})
GIT_COMMIT := $(shell git rev-parse --short=8 HEAD)
TIMESTAMP := $(shell date +%s)
PROJECT_NAME := $(notdir $(PWD))

export PROJECT_NAME

.PHONY: build graphql run stop test up

build:
	go build -v -ldflags="-X 'gqlgen-starter/cmd/build.BuildTime=$(shell date -u)' -X 'gqlgen-starter/cmd/build.BuildCommit=$(GIT_COMMIT)' -X 'gqlgen-starter/cmd/build.BuildVersion=$(VERSION)'" -o ./bin/$(PROJECT_NAME) ./cmd

clean:
	docker-compose down --remove-orphans --rmi all

graphql:
	go run github.com/99designs/gqlgen generate

run:
	air

shell:
	docker-compose -p $(PROJECT_NAME) run --rm $(PROJECT_NAME) sh

stop:
	docker-compose stop

test:
	go test ./internal/... -v -coverprofile tmp/cover.out && go tool cover -func tmp/cover.out

up:
	docker-compose up