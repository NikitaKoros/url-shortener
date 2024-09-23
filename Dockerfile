FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /url-shortener ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /url-shortener .

ENV PORT=8080
ENV DB_DSN="host=postgres user=postgres password=postgres dbname=url_shortener port=5432 sslmod=disable"

EXPOSE 8080

CMD ["./url-shortener"]