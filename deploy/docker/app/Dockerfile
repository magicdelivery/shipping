# syntax=docker/dockerfile:1
FROM golang:1.22-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/shipping_service ./cmd/web_server.go
COPY ./config ./build

FROM alpine:latest
RUN apk add curl
WORKDIR /app
COPY --from=builder /app/build/ .
CMD ["./shipping_service", "--config", "core.yaml"]