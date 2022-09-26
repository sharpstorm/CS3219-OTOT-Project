FROM golang:alpine3.16
RUN apk add make
WORKDIR /build
COPY backend ./backend
COPY makefile .
RUN cd backend
RUN touch .env
RUN make build

FROM alpine:3.16
RUN apk add libc6-compat 

WORKDIR /app
COPY --from=0 /build/dist/backend .

EXPOSE 80

ENTRYPOINT /app/backend
