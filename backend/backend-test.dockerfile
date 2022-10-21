FROM golang:alpine3.16
RUN apk add make postgresql-client
WORKDIR /backend
COPY . ./
RUN touch .env

ENTRYPOINT /bin/sh
