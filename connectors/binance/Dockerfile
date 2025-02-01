FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o binance-connector main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/binance-connector .

CMD ["./binance-connector"]