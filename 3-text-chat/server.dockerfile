FROM golang:1.17

WORKDIR /app

ADD protocol ./protocol
ADD server ./server

ADD go.* ./

RUN go build -o ./app server/server.go

CMD ["./app"]