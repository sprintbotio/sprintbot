FROM alpine:3.6

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD sprintbot /usr/local/bin/sprintbot
