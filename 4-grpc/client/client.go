package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"net"
	"strings"

	pb "4-grpc/grpc"
)

const (
	PORT = 8080
)

const HelpText = "\n" +
	"Start commands:\n" +
	"\n" +
	"help                       printing this help message\n" +
	"connect <ip_addr:port>     connect to server with specified ip and port\n" +
	"\n" +
	"Commands available after connecting to server:\n" +
	"\n" +
	"username <username>        change username\n" +
	"disconnect                 disconnect from current server\n" +
	"play                       enters game lobby. Game starts when 6 players enter lobby\n" +
	"\n" +
	"Commands available in game lobby:\n" +
	"\n" +
	"leave                      leaves lobby if game is not started yet\n" +
	"kill                       vote for kill during night if you are mafia\n" +
	"vote                       vote for kill during the day\n" +
	"\n"

var (
	Username string
	UserID   uuid.NullUUID

	Client     pb.MafiaClient
	Conn       *grpc.ClientConn
	ConnCtx    context.Context
	ConnCancel context.CancelFunc
	GameCancel context.CancelFunc

	Players   = map[uuid.UUID]*Player{}
	MyRole    pb.GameAction_Roles
	RoleToStr = map[pb.GameAction_Roles]string{
		pb.GameAction_MAFIA:   "Mafia",
		pb.GameAction_CITIZEN: "Citizen",
	}
	IsDay bool

	Actions chan *pb.GameAction
	Chat    chan string
)

var commands = map[string]func([]string) error{
	"help":       help,
	"connect":    connect,
	"username":   setUsername,
	"disconnect": disconnect,
	"play":       play,
	"leave":      leave,
	"kill":       kill,
	"vote":       vote,
}

func ProcessCommand(command string) {
	tokens := strings.Fields(command)
	if len(tokens) == 0 {
		return
	}
	tokens[0] = strings.ToLower(tokens[0])
	f, ok := commands[tokens[0]]
	if !ok {
		WriteToLog(fmt.Sprintf("Unknown command: %s\n", tokens[0]))
		return
	}
	err := f(tokens[1:])
	if err != nil {
		WriteToLog(fmt.Sprintf("Error happend on command processing: %v\n", err))
	}

}

func help(args []string) error {
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command help")
	}
	WriteToLog(HelpText)
	return nil
}

func connect(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command connect")
	}

	ip := net.ParseIP(args[0])
	if ip == nil {
		return errors.New("bad server address")
	}

	var err error
	if Conn != nil {
		Conn.Close()
	}
	addr := fmt.Sprintf("%s:%d", ip.String(), PORT)
	Conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error on connection to server: %v", err)
	}
	Client = pb.NewMafiaClient(Conn)
	ConnCtx, ConnCancel = context.WithCancel(context.Background())
	WriteToLog(fmt.Sprintf("Connected to server: %s\n", addr))
	return nil
}

func setUsername(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command username")
	}
	if Conn == nil {
		return errors.New("could not set username, you're not connected to server")
	}
	req := &pb.Username{Username: args[0]}
	if UserID.Valid {
		req.UserID = convertUUIDtoPB(UserID.UUID)
	}
	response, err := Client.SetUsername(ConnCtx, req)
	if err != nil {
		return fmt.Errorf("error on get id: %v", err)
	}
	if !UserID.Valid {
		UserID.UUID = convertPBtoUUID(response)
		UserID.Valid = true
	}
	Username = args[0]
	WriteToLog(fmt.Sprintf("Username was set to: %s\n", Username))
	return nil
}

func play(args []string) error {
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command play")
	}
	if !UserID.Valid {
		return errors.New("could not start game with empty username")
	}

	gameCtx, gameCancel := context.WithCancel(context.Background())
	GameCancel = gameCancel

	players, err := Client.GetPlayers(ConnCtx, convertUUIDtoPB(UserID.UUID))
	if err != nil {
		return fmt.Errorf("error on get player list: %v", err)
	}
	go getPlayers(gameCtx, players)

	chat, err := Client.StartChat(ConnCtx)
	if err != nil {
		return fmt.Errorf("error on starting chat: %v", err)
	}
	go startChat(gameCtx, chat)

	game, err := Client.PlayGame(ConnCtx)
	if err != nil {
		return fmt.Errorf("error on connecting to game room: %v", err)
	}
	go startGame(gameCtx, game)
	WriteToLog("You have entered the game\n")
	return nil
}

func startGame(ctx context.Context, game pb.Mafia_PlayGameClient) {
	Actions = make(chan *pb.GameAction)
	go func() {
		<-ctx.Done()
		game.Send(&pb.GameAction{
			Action: pb.GameAction_LEAVE,
			UserID: convertUUIDtoPB(UserID.UUID),
		})
		close(Actions)
		game.CloseSend()
	}()
	go func() {
		for {
			action, err := game.Recv()
			if err != nil {
				if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
					WriteToLog(fmt.Sprintf("Error on getting game action: %v\n", err))
				}
				return
			}
			processAction(action)
		}
	}()
	err := game.Send(&pb.GameAction{
		Action: pb.GameAction_ENTER,
		UserID: convertUUIDtoPB(UserID.UUID),
	})
	if err != nil {
		if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
			WriteToLog(fmt.Sprintf("Error on sending game action: %v\n", err))
		}
		return
	}
	for {
		action, more := <-Actions
		if !more {
			return
		}
		err := game.Send(action)
		if err != nil {

			if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
				WriteToLog(fmt.Sprintf("Error on sending game action: %v\n", err))
			}
			return
		}
	}
}

func switchDay() {
	if IsDay {
		WriteToLog("Night begins\n")
	} else {
		WriteToLog("Day begins\n")
	}
	IsDay = !IsDay
}

func processAction(action *pb.GameAction) {
	switch action.GetAction() {
	case pb.GameAction_ROLE:
		MyRole = action.GetRole()
		IsDay = false
		WriteToLog(fmt.Sprintf("Game started! Your role is %s\nNight begins\n", RoleToStr[MyRole]))
	case pb.GameAction_KILL:
		rawID := action.GetUserID()
		if rawID == nil {
			WriteToLog("Nobody was killed\n")
			switchDay()
			return
		}
		userID := convertPBtoUUID(rawID)
		Players[userID].status = DEAD
		if userID == UserID.UUID {
			WriteToLog("You was killed\n")
		} else {
			WriteToLog(fmt.Sprintf("%s was killed\n", Players[userID].username))
		}
		UpdateUserList()
		switchDay()
	case pb.GameAction_MAF_WIN:
		WriteToLog("Mafia win!\n")
		leave([]string{})
	case pb.GameAction_CIT_WIN:
		WriteToLog("Citizens win!\n")
		leave([]string{})
	case pb.GameAction_INFO:
		WriteToLog(action.GetData())
	}
}

type PlayerStatus uint8

const (
	ALIVE PlayerStatus = iota
	DEAD
)

type Player struct {
	username string
	status   PlayerStatus
}

func getPlayers(ctx context.Context, players pb.Mafia_GetPlayersClient) {
	Players[UserID.UUID] = &Player{
		username: Username,
		status:   ALIVE,
	}
	go func() {
		<-ctx.Done()
		players.CloseSend()
	}()
	for {
		UpdateUserList()
		change, err := players.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
				WriteToLog(fmt.Sprintf("Error on getting players change: %v\n", err))
			}
			return
		}
		userID := convertPBtoUUID(change.GetUserID())
		switch change.GetAction() {
		case pb.PlayersChange_ADD:
			Players[userID] = &Player{
				username: change.GetUsername(),
				status:   ALIVE,
			}
		case pb.PlayersChange_REMOVE:
			delete(Players, userID)
		}
	}
}

func startChat(ctx context.Context, chat pb.Mafia_StartChatClient) {
	Chat = make(chan string)
	go func() {
		<-ctx.Done()
		chat.CloseSend()
		close(Chat)
	}()
	go func() {
		err := chat.Send(&pb.ChatMessage{
			UserID: convertUUIDtoPB(UserID.UUID),
			Text:   "",
		})
		if err != nil {
			if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
				WriteToLog(fmt.Sprintf("Error on sending chat msg: %v\n", err))
			}
			return
		}
		for {
			text, more := <-Chat
			if !more {
				return
			}
			err := chat.Send(&pb.ChatMessage{
				UserID: convertUUIDtoPB(UserID.UUID),
				Text:   text,
			})
			if err != nil {
				if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
					WriteToLog(fmt.Sprintf("Error on sending chat msg: %v\n", err))
				}
				return
			}
		}
	}()

	for {
		msg, err := chat.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
				WriteToLog(fmt.Sprintf("Error on getting chat msg: %v\n", err))
			}
			return
		}
		userID := convertPBtoUUID(msg.GetUserID())
		WriteToChat(msg.Text, *Players[userID])
	}
}

func leave(args []string) error {
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command leave")
	}
	if GameCancel == nil {
		return errors.New("you can't leave you're not in a game")
	}
	GameCancel()
	GameCancel = nil
	Players = make(map[uuid.UUID]*Player)
	UpdateUserList()
	return nil
}

// Костыль - не бейте палками
func foundUser(username string) *uuid.UUID {
	for k, v := range Players {
		if v.username == username {
			return &k
		}
	}
	return nil
}

func kill(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command kill")
	}
	if MyRole != pb.GameAction_MAFIA {
		return errors.New("you're not mafia")
	}
	if Players[UserID.UUID].status == DEAD {
		return errors.New("you're dead, you can't kill")
	}
	userID := foundUser(args[0])
	if userID == nil {
		return errors.New("such user does not exists")
	} else if Players[*userID].status == DEAD {
		return errors.New("you can kill only alive users")
	}
	if IsDay {
		return errors.New("you can kill only during the night")
	}
	Actions <- &pb.GameAction{
		Action: pb.GameAction_KILL,
		Data:   stringPtr(args[0]),
	}
	WriteToLog(fmt.Sprintf("You voted to kill %s\n", args[0]))
	return nil
}

func vote(args []string) error {
	if len(args) != 1 {
		return errors.New("bad arguments for command kill")
	}
	if Players[UserID.UUID].status == DEAD {
		return errors.New("you're dead, you can't vote")
	}
	userID := foundUser(args[0])
	if userID == nil {
		return errors.New("such user does not exists")
	} else if Players[*userID].status == DEAD {
		return errors.New("you can vote for only alive users")
	}
	if !IsDay {
		return errors.New("you can vote only during the day")
	}
	Actions <- &pb.GameAction{
		Action: pb.GameAction_VOTE,
		Data:   stringPtr(args[0]),
	}
	WriteToLog(fmt.Sprintf("You voted to kill %s\n", args[0]))
	return nil
}

func disconnect(args []string) error {
	if len(args) != 0 {
		return errors.New("unnecessary arguments for command disconnect")
	}
	if Conn != nil {
		ConnCancel()
		Conn.Close()
		Conn = nil
	}
	return nil
}

func convertUUIDtoPB(id uuid.UUID) *pb.UUID {
	rawID, _ := id.MarshalBinary()
	return &pb.UUID{Value: rawID}
}

func convertPBtoUUID(id *pb.UUID) uuid.UUID {
	var userID uuid.UUID
	userID.UnmarshalBinary(id.GetValue())
	return userID
}

func stringPtr(val string) *string {
	return &val
}
