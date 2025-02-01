FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bybit-connector main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/bybit-connector .

CMD ["./bybit-connector"]