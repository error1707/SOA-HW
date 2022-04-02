FROM golang:1.17

WORKDIR /app

ADD grpc ./grpc
ADD server ./server

ADD go.* ./

RUN go build -o ./app ./server/

CMD ["./app"]