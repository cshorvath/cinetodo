# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates && adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=builder /app/server /app/server
USER app
EXPOSE 8000
ENV PORT=8000
ENTRYPOINT ["/app/server"]
