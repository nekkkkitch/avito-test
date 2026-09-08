FROM golang:1.26 AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /market_app

COPY . .

RUN go mod tidy
RUN go build -C cmd/main/ -o market_app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /market_app/cmd/main/market_app .

CMD ["./market_app"]