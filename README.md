# gqlgen-starter

This is a starter kit for building a GraphQL API using Go. It's based on my preferences for particular
technologies, so it may not suit you perfectly. At the very least, I hope it helps construct your own application.

* [gqlgen](https://github.com/99designs/gqlgen)
* Postgres and [pgx](https://github.com/JackC/pgx)
* [Zerolog](https://github.com/obsidiandynamics/zerolog)

## Getting Started

### Local Binary

1. `go mod download`
2. `touch .env` and copy contents of `.env.example` in there. Customize if you want.
3. `make run`

### Dockerized

1. `touch .env` and copy contents of `.env.example` in there. Customize if you want.
2. `make up`


Whichever method you choose to start the server, you should see log output that looks like

```
gqlgen-starter         | 2022-10-23T04:29:59Z INF internal/server.go:22 > ******************************************
gqlgen-starter         | 2022-10-23T04:29:59Z INF internal/server.go:23 >       Build Version: 0.0.1
gqlgen-starter         | 2022-10-23T04:29:59Z INF internal/server.go:24 >       Build Commit: dc7e9bd7
gqlgen-starter         | 2022-10-23T04:29:59Z INF internal/server.go:25 >       Build Time: Sun Oct 23 04:29:57 UTC 2022
gqlgen-starter         | 2022-10-23T04:29:59Z INF internal/server.go:26 > ******************************************
```

## Running Tests

It's easy! Running `make test` will run unit tests in the internal/ dir. This will also output test coverage numbers.

## Changing the GraphQL Schema

To add to the schema, create a new .graphql file in internal/graph.

Any changes to the schema should be followed by running `make graphql` to regenerate code.

## Changing the Database Schema

After changing the ent schema, the physical database schema needs to be updated too. Using the [atlas](https://atlasgo.io) 
migration engine, this code can be generated

```bash
% NAME=create_users make migration
```

and applied locally

```bash
% make migration-apply
```

## To Dos
- [ ] Add an ORM? Ent!
- [ ] Add a CI pipeline for GitHub Action and GitLab CI