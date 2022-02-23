package main

import (
	"3-voice-chat/protocol"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"os"
	"strings"
)

var (
	IDUsername = map[uuid.UUID]string{}
	UsernameID = map[string]uuid.UUID{}
	UserRoom   = map[uuid.UUID]uuid.UUID{}

	RoomID    = map[string]uuid.UUID{}
	RoomUsers = map[uuid.UUID]map[uuid.UUID]struct{}{}

	ClientsAudio = map[uuid.UUID]chan protocol.Message{}
)

func GetRoomList() string {
	keys := make([]string, len(RoomID))
	i := 0
	for k := range RoomID {
		keys[i] = k
		i++
	}
	return strings.Join(keys, ", ")
}

func ProcessMessage(msg protocol.Message, conn *net.TCPConn) protocol.Message {
	switch msg.Type {
	case protocol.CONNECT:
		username := string(msg.Data)
		_, exists := UsernameID[username]
		if exists {
			msg.Type = protocol.ERROR
			msg.Data = []byte(fmt.Sprintf("User with username: %s, already exists", username))
			return msg
		}
		id := uuid.New()
		UsernameID[username] = id
		IDUsername[id] = username
		msg.UserID.UUID = id
		msg.UserID.Valid = true
		users, _ := json.Marshal(IDUsername)
		msg.Data = users
		fmt.Printf("Set uuid: %s, for user: %s\n", id.String(), username)
		newUser := protocol.Message{
			Type:   protocol.NEW_USER,
			UserID: msg.UserID,
			Data:   []byte(username),
		}
		for _, ch := range ClientsAudio {
			ch <- newUser
		}
		_, exists = ClientsAudio[id]
		if !exists {
			ClientsAudio[id] = make(chan protocol.Message)
			go BroadcastMessage(conn, ClientsAudio[id])
		}
	case protocol.DISCONNECT:
		delete(UsernameID, IDUsername[msg.UserID.UUID])
		close(ClientsAudio[msg.UserID.UUID])
		delete(ClientsAudio, msg.UserID.UUID)
		delete(IDUsername, msg.UserID.UUID)
	case protocol.ROOM_LIST:
		msg.Data = []byte(GetRoomList())
	case protocol.ROOM_CONNECT:
		roomName := string(msg.Data)
		roomID, exists := RoomID[roomName]
		if !exists {
			roomID = uuid.New()
			RoomID[roomName] = roomID
		}
		_, exists = RoomUsers[roomID]
		if !exists {
			RoomUsers[roomID] = map[uuid.UUID]struct{}{}
		}
		RoomUsers[roomID][msg.UserID.UUID] = struct{}{}
		UserRoom[msg.UserID.UUID] = roomID
	case protocol.ROOM_DISCONNECT:
		roomID := UserRoom[msg.UserID.UUID]
		delete(UserRoom, msg.UserID.UUID)
		delete(RoomUsers[roomID], msg.UserID.UUID)
		if len(RoomUsers[roomID]) == 0 {
			delete(RoomUsers, roomID)
		}
	case protocol.TEXT:
		roomID, exists := UserRoom[msg.UserID.UUID]
		if !exists {
			msg.Type = protocol.ERROR
			msg.Data = []byte("You are not in the room!")
			return msg
		}
		for userID := range RoomUsers[roomID] {
			if msg.UserID.UUID == userID {
				continue
			}
			ClientsAudio[userID] <- msg
		}
	}
	return msg
}

func ListenCommands(conn *net.TCPConn) {
	var msg protocol.Message
	for {
		size := make([]byte, 4)
		_, err := conn.Read(size)
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}
			fmt.Fprintf(os.Stderr, "Error on recieving message: %v\n", err)
			continue
		}
		payload := make([]byte, binary.BigEndian.Uint32(size))
		_, err = conn.Read(payload)
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}
			fmt.Fprintf(os.Stderr, "Error on recieving message: %v\n", err)
			continue
		}
		if payload == nil {
			continue
		}
		msg.Decode(payload)
		resp := ProcessMessage(msg, conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not proccess command: %v\n", err)
			continue
		}
		if resp.Type == protocol.DISCONNECT {
			conn.Close()
			return
		}
		payload, err = resp.Encode()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on encoding response for command: %v\n", err)
			continue
		}
		_, err = conn.Write(payload)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Fprintf(os.Stderr, "Error on sending message: %v\n", err)
		}
	}
}

func BroadcastMessage(conn *net.TCPConn, ch chan protocol.Message) {
	var buffer []byte
	for {
		packet, opened := <-ch
		if !opened {
			fmt.Printf("Connection closed\n")
			conn.Close()
			return
		}
		buffer, _ = packet.Encode()
		_, err := conn.Write(buffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on write audio packet: %v\n", err)
		}
	}
}

func main() {
	addrRaw := fmt.Sprintf("0.0.0.0:%d", protocol.PORT)
	addrTCP, _ := net.ResolveTCPAddr("tcp", addrRaw)
	listenerTCP, err := net.ListenTCP("tcp", addrTCP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start server: %v\n", err)
		return
	}
	fmt.Println("Server started")

	for {
		clientConn, err := listenerTCP.AcceptTCP()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on client connection: %v\n", err)
		}
		go ListenCommands(clientConn)
	}
}
