include .env
export $(shell sed 's/=.*//' .env)

versionFile := version.txt
VERSION := $(shell cat ${versionFile})
GIT_COMMIT := $(shell git rev-parse --short=8 HEAD)
TIMESTAMP := $(shell date +%s)
PROJECT_NAME := $(notdir $(PWD))
HOST_DB_URL := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

export PROJECT_NAME

.PHONY: build clean ent ent-init graphql migration migration-db migration-hash migration-lint migrate-apply-local run stop test up

build:
	go build -v -ldflags="-X 'gqlgen-starter/cmd/build.BuildTime=$(shell date -u)' -X 'gqlgen-starter/cmd/build.BuildCommit=$(GIT_COMMIT)' -X 'gqlgen-starter/cmd/build.BuildVersion=$(VERSION)'" -o ./bin/$(PROJECT_NAME) ./cmd

clean:
	docker-compose down --remove-orphans --rmi all

# Generates supporting ent code whenever the ent schema changes
ent:
	go generate ./internal/ent

# Creates a new model in the ent schema
ent-init:
	go run -mod=mod entgo.io/ent/cmd/ent init --target ./internal/ent $(MODEL)

# Generate gql related code after the GraphQL schema changes
graphql:
	go run github.com/99designs/gqlgen generate

# Creates a new migration file in db/migrations. Use like \
> NAME=create_users make migration
migration:
	go run -mod=mod internal/ent/migrate/main.go $(NAME)

# Starts a containerized Postgres database for the atlas engine. See https://atlasgo.io/concepts/dev-database
migration-db:
	docker run --name atlas-migration --rm -p 5678:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=test -d postgres

# When the hash sum is invalid because you manually updated a file after the hash sum was generated, this command \
will force Atlas to recompute the value
migration-hash:
	atlas migrate hash --dir="file://$(PWD)/db/migrations"

# If the "dev database" is running from the migration-db target, then this will analyze the code in db/migrations for \
potentially dangerous changes
migration-lint:
	atlas migrate lint --dev-url="postgres://postgres:postgres@localhost:5678/test?sslmode=disable" --dir="file://$(PWD)/db/migrations" --latest=1

# Applies the latest migrations in the local development database
migrate-apply-local:
	atlas migrate apply --dir="file://$(PWD)/db/migrations" --url=$(HOST_DB_URL)

# Run the API locally in a hot reload environment
run:
	air

# Open a bash shell into the API container
shell:
	docker-compose -p $(PROJECT_NAME) run --rm $(PROJECT_NAME) sh

stop:
	docker-compose stop

test:
	go test ./internal/... -v -coverprofile tmp/cover.out && go tool cover -func tmp/cover.out

up:
	docker-compose up