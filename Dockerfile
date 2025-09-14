FROM golang:1.24.7-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main main.go

FROM debian:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/migrations ./migrations

CMD ["./main"]
