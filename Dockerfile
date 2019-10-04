FROM golang:1.13

LABEL maintainer="andygrunwald@gmail.com"

# Default HTTP Webserver port
EXPOSE 8080

# Default software buzzer emulation port
EXPOSE 8181

WORKDIR /go/src/app
COPY . .

RUN go build -i -o twb-websocket

ENTRYPOINT ["./twb-websocket"]