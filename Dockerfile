FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o log-analyzer .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/log-analyzer .
COPY config/ config/
ENTRYPOINT ["./log-analyzer"]
