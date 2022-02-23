package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"3-voice-chat/protocol"
)

var (
	Username  string
	UserID    uuid.NullUUID
	Users     map[uuid.UUID]string
	RoomName  string
	Connected bool

	Disconnect func()
	Send       chan protocol.Message
)

func ProcessControl(msg protocol.Message) error {
	switch msg.Type {
	case protocol.CONNECT:
		UserID = msg.UserID
		fmt.Printf("Got uuid: %s\n", UserID.UUID.String())
		err := json.Unmarshal(msg.Data, &Users)
		if err != nil {
			fmt.Printf("Error on getting users")
		}
	case protocol.NEW_USER:
		userID := msg.UserID.UUID
		username := string(msg.Data)
		Users[userID] = username
		fmt.Printf("New user: %s\n", username)
	case protocol.ROOM_LIST:
		fmt.Printf("Available rooms: %s\n", msg.Data)
	case protocol.ROOM_CONNECT:
		RoomName = string(msg.Data)
		fmt.Printf("Connected to room: %s\n", RoomName)
	case protocol.ROOM_DISCONNECT:
		fmt.Printf("Disconnected from room: %s\n", RoomName)
		RoomName = ""
	case protocol.TEXT:
		if msg.UserID.UUID == UserID.UUID {
			return nil
		}
		username, exists := Users[msg.UserID.UUID]
		if !exists {
			username = "unknown user"
		}
		text := string(msg.Data)
		fmt.Printf("%s >>> %s\n", username, text)
	case protocol.ERROR:
		fmt.Printf("Got error from server: %s\n", string(msg.Data))
	}
	return nil
}

func ProcessHelp(args []string) error {
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command help")
	}
	fmt.Print(
		"\n",
		"Start commands:\n",
		"\n",
		"help                       printing this help message\n",
		"username <username>        change username\n",
		"connect <ip_addr:port>     connect to server with specified ip and port\n",
		"\n",
		"Commands available after connecting to server:\n",
		"\n",
		"disconnect                 disconnect from current server\n",
		"room list                  get list of available rooms\n",
		"room connect <room_name>   connect to available room, or create new with specified name\n",
		"room disconnect            disconnect from current room\n",
		"\n",
		"Commands available in room:\n",
		"\n",
		"text <message>             send message to room\n",
		"\n",
	)
	return nil
}

func ProcessUsername(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command username")
	}
	Username = args[0]
	fmt.Printf("Username was set to: %s\n", Username)
	return nil
}

func ProcessConnect(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command connect")
	}
	if Connected {
		return errors.New("you're already connected")
	}
	if Username == "" {
		return errors.New("username not set")
	}

	ip := net.ParseIP(args[0])
	if ip == nil {
		return errors.New("bad server address")
	}

	var conn *net.TCPConn
	addrTCP := &net.TCPAddr{
		IP:   ip,
		Port: protocol.PORT,
	}
	var err error
	for i := 0; i < 3; i++ {
		conn, err = net.DialTCP("tcp", nil, addrTCP)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(1+i) * time.Second)
	}
	if err != nil {
		return err
	}
	Connected = true
	Disconnect = func() {
		conn.Close()
	}

	Send = make(chan protocol.Message)
	ProcessMessaging(conn)
	if err != nil {
		return fmt.Errorf("could not get User ID: %v", err)
	}

	return nil
}

func ProcessMessaging(conn *net.TCPConn) {
	go func() {
		for msg := range Send {
			payload, err := msg.Encode()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error on encoding message: %v\n", err)
				continue
			}
			_, err = conn.Write(payload)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error on sending message: %v\n", err)
				continue
			}
		}
	}()

	go func() {
		var msg protocol.Message
		for {
			size := make([]byte, 4)
			_, err := conn.Read(size)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}
				fmt.Fprintf(os.Stderr, "error on receiving message: %v\n", err)
				continue
			}
			payload := make([]byte, binary.BigEndian.Uint32(size))
			_, err = conn.Read(payload)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}
				fmt.Fprintf(os.Stderr, "error on receiving message: %v\n", err)
				continue
			}
			err = msg.Decode(payload)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error on decoding message: %v\n", err)
				continue
			}
			err = ProcessControl(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error on processing msg: %v\n", err)
				continue
			}
		}
	}()

	// Initial: getting User ID
	msg := protocol.Message{
		Type: protocol.CONNECT,
		Data: []byte(Username),
	}
	Send <- msg

	// Send msg to server

}

func ProcessDisconnect(args []string) error {
	if !Connected {
		return errors.New("cannot disconnect: you're not connected to server")
	}
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command disconnect")
	}
	msg := protocol.Message{
		Type:   protocol.DISCONNECT,
		UserID: UserID,
	}
	Send <- msg
	Connected = false
	Disconnect()
	UserID = uuid.NullUUID{}
	RoomName = ""
	fmt.Printf("Disconnected\n")
	return nil
}

func ProcessRoom(args []string) error {
	if !Connected {
		return errors.New("cannot interact with rooms: you're not connected to server")
	}
	if len(args) == 0 {
		return errors.New("subcommand for room command not specified")
	}
	switch args[0] {
	case "list":
		if len(args) != 1 {
			return errors.New("unnecessary args for command room list")
		}
		msg := protocol.Message{
			Type:   protocol.ROOM_LIST,
			UserID: UserID,
		}
		Send <- msg
	case "connect":
		if len(args) != 2 {
			return errors.New("bad args for command room connect")
		}
		msg := protocol.Message{
			Type:   protocol.ROOM_CONNECT,
			UserID: UserID,
			Data:   []byte(args[1]),
		}
		Send <- msg
	case "disconnect":
		if len(args) != 1 {
			return errors.New("unnecessary args for command room disconnect")
		}
		if RoomName == "" {
			return errors.New("you're not in room")
		}
		msg := protocol.Message{
			Type:   protocol.ROOM_DISCONNECT,
			UserID: UserID,
		}
		Send <- msg
	}
	return nil
}

func ProcessText(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command text")
	}
	Send <- protocol.Message{
		Type:   protocol.TEXT,
		UserID: UserID,
		Data:   []byte(strings.Join(args, " ")),
	}
	return nil
}

var commands = map[string]func([]string) error{
	"help":       ProcessHelp,
	"username":   ProcessUsername,
	"connect":    ProcessConnect,
	"disconnect": ProcessDisconnect,
	"room":       ProcessRoom,

	"text": ProcessText,
}

func main() {
	fmt.Println("Welcome to VoiceChat!")
	fmt.Println("Set username before connecting to server")
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
