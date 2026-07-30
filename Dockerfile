FROM golang:1.25 AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o shortener cmd/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/shortener .

EXPOSE 8080

ENTRYPOINT ["./shortener"]
