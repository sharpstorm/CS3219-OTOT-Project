FROM golang:alpine3.16
RUN apk add make postgresql-client
WORKDIR /backend
COPY . ./
RUN rm .env
RUN touch .env

ENTRYPOINT /bin/sh
