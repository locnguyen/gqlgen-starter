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
2022/07/17 07:43:56 ******************************************
2022/07/17 07:43:56     Build Commit: dc0b1a11
2022/07/17 07:43:56     Build Time: Sun Jul 17 14:43:55 UTC 2022
2022/07/17 07:43:56     Environment Variables: {
  "ServerPort": "9000"
}
2022/07/17 07:43:56 ******************************************
2022/07/17 07:43:56 connect to http://localhost:9000/ for GraphQL playground
```

## Running Tests

It's easy! Running `make test` will run unit tests in the internal/ dir. This will also output test coverage numbers.

## Changing the GraphQL Schema

To add to the schema, create a new .graphql file in internal/graph.

Any changes to the schema should be followed by running `make graphql` to regenerate code.

## To Dos
- [ ] Add an ORM? Ent!
- [ ] Add a CI pipeline for GitHub Action and GitLab CI