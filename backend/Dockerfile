FROM golang:1.26 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -C cmd/main/ -o app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cmd/main/app .
COPY --from=builder /app/migrations ./migrations

CMD ["./app"]