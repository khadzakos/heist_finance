FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o controller cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/controller .
COPY configs/config.yaml /root/configs/
CMD ["./controller"]