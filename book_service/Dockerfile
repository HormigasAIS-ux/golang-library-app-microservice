FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o book_service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/book_service .
