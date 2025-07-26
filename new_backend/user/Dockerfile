FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o main ./cmd

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates curl

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["./main"]
