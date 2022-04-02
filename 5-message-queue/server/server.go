package main

import (
	pb "5-message-queue/grpc"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	PORT = 8080
)

var (
	NextID      uint64 = 0
	UserResults        = map[uint64]chan *pb.Result{}
	Channel     *amqp.Channel
)

func main() {
	setupQueue()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", PORT))
	if err != nil {
		log.Fatalf("Error on listening addr: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterWikiPathLengthMeterServer(s, &Server{})
	log.Printf("Starting server\n")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error on serving: %v\n", err)
	}
}

func setupQueue() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Error on connecting to rabbitmq: %v\n", err)
	}
	//defer conn.Close()
	Channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Error on opening a channel: %v\n", err)
	}
	//defer Channel.Close()

	_, err = Channel.QueueDeclare(
		"requests", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Error on declaring a queue: %v\n", err)
	}
	Responses, err := Channel.QueueDeclare(
		"responses", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("Error on declaring a queue: %v\n", err)
	}
	go func() {
		msgs, err := Channel.Consume(
			Responses.Name, // queue
			"",             // consumer
			true,           // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		if err != nil {
			log.Fatalf("Error on registering a consumer: %v\n", err)
		}
		for d := range msgs {
			var msg Result
			json.Unmarshal(d.Body, &msg)
			log.Printf(
				"Got results: \n"+
					"URL1: %s\n"+
					"URL2: %s\n"+
					"Fom user with id %d\n",
				msg.URL1, msg.URL2, msg.UserID)
			UserResults[msg.UserID] <- &pb.Result{
				URL1: msg.URL1,
				URL2: msg.URL2,
				Path: msg.Path,
			}
		}
	}()
}

type Result struct {
	UserID uint64   `json:"UserID"`
	URL1   string   `json:"URL1"`
	URL2   string   `json:"URL2"`
	Path   []string `json:"Path"`
}

type Request struct {
	UserID uint64 `json:"UserID"`
	URL1   string `json:"URL1"`
	URL2   string `json:"URL2"`
}

type Server struct {
	pb.UnimplementedWikiPathLengthMeterServer
}

func (s *Server) GetID(ctx context.Context, _ *pb.Empty) (*pb.UserID, error) {
	NextID++
	return &pb.UserID{ID: NextID}, nil
}

func (s *Server) MeasurePath(ctx context.Context, req *pb.PathRequest) (*pb.Empty, error) {
	log.Printf(
		"Got request: \n"+
			"URL1: %s\n"+
			"URL2: %s\n"+
			"From user with id %d\n",
		req.URL1, req.URL2, req.User.GetID())
	task := Request{
		UserID: req.User.GetID(),
		URL1:   req.URL1,
		URL2:   req.URL2,
	}
	rawTask, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	err = Channel.Publish("", // exchange
		"requests", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        rawTask,
		})
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *Server) GetResults(userID *pb.UserID, results pb.WikiPathLengthMeter_GetResultsServer) error {
	resChan := make(chan *pb.Result)
	UserResults[userID.GetID()] = resChan
	for {
		resp, more := <-resChan
		if !more {
			return nil
		}
		err := results.Send(resp)
		if err != nil {
			return err
		}
	}
}
