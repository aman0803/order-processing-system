FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/order-service .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./order-service"]