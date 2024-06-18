FROM golang:alpine as builder
LABEL org.opencontainers.image.source https://github.com/oakestra/plugin-kubernetes

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o oakestra-agent
RUN ls -l

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/oakestra-agent .

CMD ["./oakestra-agent"]