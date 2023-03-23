FROM golang:1.20

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o /usr/local/bin/service ./cmd/srv/...

EXPOSE 80

ENTRYPOINT service
