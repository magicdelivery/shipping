# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/shipping_service ./cmd/shipping_service.go
COPY ./config ./build

FROM alpine:latest
RUN apk add curl \
    && adduser -D appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/build/ .
CMD ["./shipping_service", "--config", "core.yaml"]