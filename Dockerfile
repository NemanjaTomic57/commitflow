FROM golang:1.25-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /out/producer ./cmd/producer/producer.go
RUN go build -o /out/consumer ./cmd/consumer/consumer.go

FROM debian:trixie-slim AS app

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /out/producer /usr/local/bin/producer
COPY --from=builder /out/consumer /usr/local/bin/consumer
COPY migrations migrations
