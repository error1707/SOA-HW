package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"net/url"
	"os"
	"strings"

	pb "5-message-queue/grpc"
)

const (
	PORT = 8080
)

var (
	ID uint64

	Conn   *grpc.ClientConn
	Client pb.WikiPathLengthMeterClient

	ConnCtx context.Context
)

func help(args []string) error {
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command help")
	}
	fmt.Print(
		"\n",
		"Start commands:\n",
		"\n",
		"help                       printing this help message\n",
		"connect <ip_addr>          connect to server with specified ip\n",
		"\n",
		"Commands available after connecting to server:\n",
		"\n",
		"disconnect                 disconnect from current server\n",
		"find <URL1> <URL2>         find path from URL1 to URL2\n",
		"\n",
	)
	return nil
}

func onFail() {
	Conn.Close()
	Conn = nil
}

func connect(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command connect")
	}
	if Conn != nil {
		return errors.New("you're already connected")
	}

	ip := net.ParseIP(args[0])
	if ip == nil {
		return errors.New("bad server address")
	}
	addr := fmt.Sprintf("%s:%d", ip.String(), PORT)

	var err error
	Conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error on connection to server: %v", err)
	}
	Client = pb.NewWikiPathLengthMeterClient(Conn)
	ConnCtx = context.Background()

	resp, err := Client.GetID(ConnCtx, &pb.Empty{})
	if err != nil {
		onFail()
		return err
	}
	ID = resp.GetID()
	go getResults()
	fmt.Printf("Connected to %s\n", args[0])
	return nil
}

func getResults() {
	results, err := Client.GetResults(ConnCtx, &pb.UserID{ID: ID})
	if err != nil {
		fmt.Printf("error on creating result stream\n")
		onFail()
		return
	}
	for {
		result, err := results.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("error on getting result: %v\n", err)
			}
			return
		}
		printResult(result)
	}
}

func printResult(result *pb.Result) {
	fmt.Printf("Distance is %d between\n%s\nand\n%s\n", len(result.Path), result.URL1, result.URL2)
	fmt.Printf("And path:\n")
	for _, link := range result.Path {
		fmt.Printf("%s\n", link)
	}
}

func disconnect(args []string) error {
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command disconnect")
	}
	Conn.Close()
	Conn = nil
	fmt.Printf("Disconnected from server\n")
	return nil
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func find(args []string) error {
	if len(args) != 2 {
		return errors.New("bad arguments for command measure")
	}
	for _, urlLink := range args {
		if !isUrl(urlLink) {
			return fmt.Errorf("bad url: %s", urlLink)
		}
	}
	_, err := Client.MeasurePath(ConnCtx, &pb.PathRequest{
		User: &pb.UserID{ID: ID},
		URL1: args[0],
		URL2: args[1],
	})
	return err
}

var commands = map[string]func([]string) error{
	"help":       help,
	"connect":    connect,
	"disconnect": disconnect,
	"find":       find,
}

func main() {
	fmt.Println("Welcome to WikiPathFinder!")
	fmt.Println("Write 'help' to list available commands")

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		command := input.Text()
		if err := input.Err(); err != nil {
			fmt.Printf("Error happend on parsing command: %v\n", err)
			continue
		}
		tokens := strings.Fields(command)
		if len(tokens) == 0 {
			continue
		}
		tokens[0] = strings.ToLower(tokens[0])
		f, ok := commands[tokens[0]]
		if !ok {
			fmt.Printf("Unknown command: %s\n", tokens[0])
			continue
		}
		err := f(tokens[1:])
		if err != nil {
			fmt.Printf("Error happend on command processing: %v\n", err)
			continue
		}
	}
}
