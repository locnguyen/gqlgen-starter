FROM golang:alpine AS env
RUN apk --no-cache add ca-certificates make curl build-base git

FROM env AS base
ARG PROJECT_NAME=${PROJECT_NAME:-nodummy}
WORKDIR /home/${PROJECT_NAME}
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

CMD air