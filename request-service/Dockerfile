FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o request-service cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/request-service .

CMD ["./request-service"]