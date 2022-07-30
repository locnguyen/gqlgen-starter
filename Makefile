include .env
export $(shell sed 's/=.*//' .env)

GIT_COMMIT := $(shell git rev-parse --short=8 HEAD)
TIMESTAMP := $(shell date +%s)
PROJECT_NAME = $(notdir $(PWD))

.PHONY: build graphql run

build:
	go build -v -ldflags="-X 'main.BuildTime=$(shell date -u)' -X 'main.BuildCommit=$(GIT_COMMIT)'" -o ./bin/$(PROJECT_NAME) ./cmd

graphql:
	go run github.com/99designs/gqlgen generate

run:
	air