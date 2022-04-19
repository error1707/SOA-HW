FROM golang:1.17

WORKDIR /app

ADD server ./server
ADD views ./views

ADD go.* ./

RUN go build -o ./app ./server/

CMD ["./app"]