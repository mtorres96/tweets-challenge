# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/server /app/server
# Defaults de rate limit (podés override con `-e` en runtime)
ENV RATE_LIMIT_ENABLED=true \
    RATE_LIMIT_WINDOW_SEC=60 \
    RATE_LIMIT_MAX_TWEETS=20
# El puerto real lo define la app leyendo PORT; EXPOSE es solo documental
EXPOSE 8080
ENTRYPOINT ["/app/server"]
