# Используем официальный образ Go как базовый
FROM golang:latest as builder


WORKDIR /app
ENV GOOS=linux
ENV CGO_ENABLED=1
COPY . .
RUN go mod tidy
RUN  go build -C ./cmd/api  -o ../../bin


FROM debian:bookworm
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/bin .
COPY --from=builder /app/.env .
EXPOSE 8080


CMD ["./bin"]

