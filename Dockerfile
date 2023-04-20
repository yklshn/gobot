# FROM golang:1.19 AS builder

# COPY . /GorBot
# WORKDIR /GorBot

# RUN  go build -o ./bin/main ./cmd

# FROM alpine
# WORKDIR /

# COPY --from=builder GorBot/bin/ /app

# CMD ./app/main

FROM golang:1.18.5-alpine3.15 AS builder

COPY . /Bot
WORKDIR /Bot

RUN  go build -o ./bin/main ./cmd

FROM alpine:3.15
WORKDIR /root/

COPY --from=0 /Bot/bin/main .
COPY --from=0 /Bot/config ./config

EXPOSE 3333

CMD ["./main"]