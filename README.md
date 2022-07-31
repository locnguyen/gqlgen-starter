# gqlgen-starter

This is a simple starter kit for anyone wanting to build a GraphQL API using [gqlgen](https://github.com/99designs/gqlgen).

## Getting Started

### Local Binary

1. `go mod download`
2. `touch .env`
3. `make run`

### Dockerized

1. `touch .env`
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
- [ ] Add a database connection
- [ ] Add an ORM?
- [ ] Use structured logging
- [x] Example unit tests
- [x] Containerize this
- [ ] Add a CI pipeline for GitHub Action and GitLab CI