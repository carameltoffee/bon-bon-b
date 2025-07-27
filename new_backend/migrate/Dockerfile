FROM golang:1.24-alpine

RUN apk add --no-cache git
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app
COPY . /app
