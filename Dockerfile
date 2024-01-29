FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY ./ ./

RUN go build -o server_app

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server_app .

ENTRYPOINT ["./server_app"]
