# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
RUN apk add --no-cache build-base ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# gorm.io/driver/sqlite necesita cgo
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /app/server /app/server

# Defaults editables por env
ENV RATE_LIMIT_ENABLED=true \
    RATE_LIMIT_WINDOW_SEC=60 \
    RATE_LIMIT_MAX_TWEETS=20
EXPOSE 8080
ENTRYPOINT ["/app/server"]
