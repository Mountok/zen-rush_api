# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o zenrush-backend ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/zenrush-backend ./zenrush-backend
EXPOSE 8080
CMD ["./zenrush-backend"] 