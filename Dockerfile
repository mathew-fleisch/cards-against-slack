FROM golang:alpine AS builder

LABEL maintainer="Mathew Fleisch <mathew.fleisch@gmail.com>"

ENV SLACK_TOKEN=

RUN mkdir /cards
WORKDIR /cards
COPY . .
RUN apk update && apk add --no-cache git bash
RUN go install -v ./...
RUN go build bot.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /cards
COPY . .
COPY --from=builder /cards/bot /usr/local/bin/bot
CMD ["bot","files/questions.txt","files/answers.txt","files/triggers.txt"]
