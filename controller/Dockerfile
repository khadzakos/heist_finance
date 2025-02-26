FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o controller cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/controller .
COPY --from=builder /app/configs/ ./configs/

CMD ["./controller"]