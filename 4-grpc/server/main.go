package main

import (
	pb "4-grpc/grpc"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	PORT = 8080
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", PORT))
	if err != nil {
		log.Fatalf("Error on listening addr: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMafiaServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error on serving: %v", err)
	}
}
