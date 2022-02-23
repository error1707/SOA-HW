FROM golang:1.17

WORKDIR /app

ADD protocol ./protocol
ADD client ./client

ADD go.* ./

RUN go build -o ./app client/client.go

CMD ["./app"]
#CMD ["/bin/sh"]