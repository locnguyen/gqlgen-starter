include .env
export $(shell sed 's/=.*//' .env)

versionFile := version.txt
VERSION := $(shell cat ${versionFile})
GIT_COMMIT := $(shell git rev-parse --short=8 HEAD)
TIMESTAMP := $(shell date +%s)
PROJECT_NAME := $(notdir $(PWD))
HOST_DB_URL := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
HOST_TEST_DB_URL := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/test?sslmode=disable"

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
	go run -mod=mod entgo.io/ent/cmd/ent init --target ./internal/ent/schema $(MODEL)

# Generate gql related code after the GraphQL schema changes
graphql:
	go run github.com/99designs/gqlgen generate

lint:
	staticcheck ./internal/...

# Creates a new migration file in db/migrations. Atlas will first inspect the test database and perform a diff
migration:
	go run -mod=mod internal/ent/migrate/main.go $(NAME)

# When the hash sum is invalid because you manually updated a file after the hash sum was generated, this command \
will force Atlas to recompute the value
migration-hash:
	atlas migrate hash --dir="file://db/migrations"

# If the "dev database" is running from the migration-db target, then this will analyze the code in db/migrations for \
potentially dangerous changes
migration-lint:
	atlas migrate lint --dev-url=$(TEST_DB_URL) --dir="file://$(PWD)/db/migrations" --latest=1

# Applies the latest migrations in the local development database
migrate-apply-local:
	atlas migrate apply --dir="file://db/migrations" --url=$(HOST_DB_URL)

# Applies the latest migrations in the local development database
migrate-apply-test:
	atlas migrate apply --dir="file://db/migrations" --url=$(TEST_DB_URL)

# Nuke the database schema in the local database
migrate-clean-local:
	atlas schema clean --url=$(HOST_DB_URL)

# Nuke the database schema in the atlas test database
migrate-clean-test:
	atlas schema clean --url=$(TEST_DB_URL)

# Run the API locally in a hot reload environment
run:
	air

# Open a bash shell into the API container
shell:
	docker-compose -p $(PROJECT_NAME) run --rm $(PROJECT_NAME) sh

stop:
	docker-compose stop

test:
	go test ./... -coverprofile tmp/cover.pre && cat tmp/cover.pre | (grep -v ./internal/ent | grep -v .internal/gql) > tmp/cover.final && go tool cover -func tmp/cover.final

up:
	docker-compose up