FROM golang:1.18-alpine3.15 AS builder

COPY . /github.com/mbakumenkov/go-pocket-bot
WORKDIR /github.com/mbakumenkov/go-pocket-bot

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=0 /github.com/mbakumenkov/go-pocket-bot/bin/bot .
COPY --from=0 /github.com/mbakumenkov/go-pocket-bot/configs configs/

EXPOSE 80

CMD ["./bot"]
