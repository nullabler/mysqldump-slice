FROM golang:1.19.4-alpine3.17

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

RUN apk update && apk add mysql-client

WORKDIR /app
COPY . /app

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o cmd/main cmd/main.go" --command="./cmd/main"
