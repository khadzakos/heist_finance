FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o binance-connector main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/binance-connector .

COPY --from=builder /app/config ./config

CMD ["./binance-connector"]