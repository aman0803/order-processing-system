FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/inventory-service .
COPY --from=builder /app/.env .
CMD ["./inventory-service"]