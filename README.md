# gqlgen-starter

This is a simple starter kit for anyone wanting to build a GraphQL API using [gqlgen](https://github.com/99designs/gqlgen).

## Getting Started

1. `go mod download`
2. `touch .env`
3. `make run`

If everything is peachy then you should see log output that looks like

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

It's easy! It should also output some test coverage numbers.

`make test`

## To Dos
- [ ] Add a database connection
- [ ] Add an ORM?
- [ ] Use structured logging
- [x] Example unit tests
- [ ] Containerize this
- [ ] Add a CI pipeline for GitHub Action and GitLab CI