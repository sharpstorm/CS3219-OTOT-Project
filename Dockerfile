FROM alpine:3.16
RUN apk add libc6-compat 

WORKDIR /app
COPY dist/backend .

EXPOSE 80

ENTRYPOINT /app/backend
