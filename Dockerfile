FROM golang:alpine3.16
RUN apk add make
WORKDIR /build
COPY backend ./backend
COPY makefile .
RUN touch backend/.env
RUN make build

FROM node:alpine as builder
WORKDIR /fe

COPY frontend/package.json frontend/package-lock.json ./frontend/
RUN cd frontend && npm i

COPY frontend/src frontend/src
COPY frontend/public frontend/public
RUN cd frontend && npm run-script build

FROM alpine:3.16
RUN apk add libc6-compat 

WORKDIR /app
COPY --from=0 /build/dist/backend .
RUN mkdir static
COPY --from=1 /fe/frontend/build static

EXPOSE 80

ENTRYPOINT /app/backend
