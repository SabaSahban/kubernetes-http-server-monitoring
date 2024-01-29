FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY ./ ./

RUN go build -o weather_app

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/weather_app .

ENTRYPOINT ["./weather_app"]