# gqlgen-starter

This is a starter kit for building a GraphQL API using Go. It's based on my preferences for particular
technologies, so it may not suit you perfectly. At the very least, I hope it helps construct your own application.

* [gqlgen](https://github.com/99designs/gqlgen)
* Postgres and [pgx](https://github.com/JackC/pgx)
* [Zerolog](https://github.com/obsidiandynamics/zerolog)
* [Ent](https://entgo.io) as the ORM
* [nats](https://nats.io/) for pub/sub messaging

## Getting Started

### Local Binary

Install the external dependencies first and ensure they are running for the binary to connect to. Then

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

## Changing the Ent and/or Database Schema

After changing the ent schema, you'll need to perform a number of manual steps (you may also need to regenerate the 
GraphQL schema).

1. `make ent` to generate updated ORM client code
2. `make migrate-apply-test` to ensure that the test DB has the latest migrations applied
3. `NAME=add_column_to_some_table make migration` which performs a diff between the test DB schema and the current ent schema to generate the SQL in a new migration file in db/migrations
4. `make migrate-apply-test` to apply the new migration file against the test DB
5. `make migrate-apply-local` to apply the new migration file against the local app DB (assuming step 4)


## To Dos
- [ ] Add a CI pipeline for GitHub Action and GitLab CI