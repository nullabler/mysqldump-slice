FROM golang:1.17.2-alpine3.14

RUN apk update && apk add mysql-client git ffmpeg ca-certificates && update-ca-certificates

WORKDIR /app
COPY . /app

RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon --build="go build -o cmd/main cmd/main.go" --command="./cmd/main"
