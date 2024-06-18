FROM golang:alpine as builder
LABEL org.opencontainers.image.source https://github.com/oakestra/plugin-kubernetes

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o netmanager

FROM debian:buster
WORKDIR /app

COPY --from=builder /app/netmanager .