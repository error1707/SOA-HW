package main

import (
	pb "4-grpc/grpc"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math/rand"
	"time"
)

const (
	MAFIA_COUNT   = 2
	CITIZEN_COUNT = 4
)

var (
	PlayerID     = map[string]uuid.UUID{}
	IDPlayer     = map[uuid.UUID]*Player{}
	InGame       = map[uuid.UUID]struct{}{}
	GameStarted  = false
	IsDay        bool
	MafiaVotes   = map[uuid.UUID]struct{}{}
	MafiaVoters  = map[uuid.UUID]uuid.UUID{}
	AllVotes     = map[uuid.UUID]int{}
	AllVoters    = map[uuid.UUID]uuid.UUID{}
	MafiaAlive   = MAFIA_COUNT
	CitizenAlive = CITIZEN_COUNT
)

func VoteAccept() int {
	return (MafiaAlive+CitizenAlive)/2 + (MafiaAlive+CitizenAlive)%2
}

type PlayerStatus uint8

const (
	ALIVE PlayerStatus = iota
	DEAD
)

type Player struct {
	username string
	status   PlayerStatus
	role     pb.GameAction_Roles

	actions chan *pb.GameAction
	players chan *pb.PlayersChange
	chat    chan *pb.ChatMessage
}

func (user *Player) SendPlayersUpdate(id uuid.UUID, action pb.PlayersChange_ChangeType) {
	msg := &pb.PlayersChange{
		Action: action,
		UserID: convertUUIDtoPB(id),
	}
	if action == pb.PlayersChange_ADD {
		msg.Username = &IDPlayer[id].username
	}
	user.players <- msg
}

func (user *Player) SendAllPlayers() {
	for player := range InGame {
		user.SendPlayersUpdate(player, pb.PlayersChange_ADD)
	}
}

func SendPlayersUpdateToAll(id uuid.UUID, action pb.PlayersChange_ChangeType) {
	for player := range InGame {
		IDPlayer[player].SendPlayersUpdate(id, action)
	}
}

type Server struct {
	pb.UnimplementedMafiaServer
}

func (s *Server) SetUsername(ctx context.Context, req *pb.Username) (*pb.UUID, error) {
	reqUUID := req.GetUserID()
	id, exists := PlayerID[req.Username]
	if exists {
		if reqUUID == nil {
			return nil, errors.New("specified username already busy")
		}
		var reqId uuid.UUID
		err := reqId.UnmarshalBinary(reqUUID.GetValue())
		if err != nil {
			return nil, fmt.Errorf("bad uuid: %v", err)
		}
		idRaw, _ := id.MarshalBinary()
		reqIdRaw, _ := reqId.MarshalBinary()
		if bytes.Compare(idRaw, reqIdRaw) != 0 {
			return nil, errors.New("specified username already busy")
		}
		return &pb.UUID{Value: idRaw}, nil
	}
	var userID uuid.UUID
	if reqUUID == nil {
		userID = uuid.New()
		IDPlayer[userID] = &Player{}
	} else {
		userID.UnmarshalBinary(reqUUID.GetValue())
		username := IDPlayer[userID].username
		delete(PlayerID, username)
	}
	PlayerID[req.Username] = userID
	IDPlayer[userID].username = req.Username
	log.Printf("Username was set: %s --> %s", req.Username, userID.String())
	userIDRaw, _ := userID.MarshalBinary()
	return &pb.UUID{Value: userIDRaw}, nil
}

func (s *Server) PlayGame(server pb.Mafia_PlayGameServer) error {
	action, err := server.Recv()
	if err != nil {
		if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
			log.Printf("Error on getting action: %v", err)
		}
		return err
	}
	userID := convertPBtoUUID(action.GetUserID())
	actions := make(chan *pb.GameAction)
	IDPlayer[userID].actions = actions
	go func() {
		for {
			action, more := <-actions
			if !more {
				return
			}
			err := server.Send(action)
			if err != nil {
				if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
					log.Printf("Error on getting action: %v", err)
				}
				return
			}
		}
	}()
	processAction(action, userID)
	for {
		action, err := server.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
				log.Printf("Error on getting action: %v", err)
			}
			return err
		}
		processAction(action, userID)
	}
}

func processAction(action *pb.GameAction, userID uuid.UUID) {
	switch action.Action {
	case pb.GameAction_ENTER:
		if GameStarted {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("Game already started, you can't enter the game\n"),
			}
			return
		}
		SendPlayersUpdateToAll(userID, pb.PlayersChange_ADD)
		IDPlayer[userID].SendAllPlayers()
		IDPlayer[userID].status = ALIVE
		InGame[userID] = struct{}{}
		log.Printf("Player %s has entered the game", IDPlayer[userID].username)
		if len(InGame) == MAFIA_COUNT+CITIZEN_COUNT {
			GameStarted = true
			sendRoles()
			IsDay = false
		}
	case pb.GameAction_LEAVE:
		if GameStarted {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("Game already started, you can't leave the game\n"),
			}
			return
		}
		SendPlayersUpdateToAll(userID, pb.PlayersChange_REMOVE)
		log.Printf("Player %s has left the game", IDPlayer[userID].username)
		delete(InGame, userID)
	case pb.GameAction_KILL:
		_, exists := MafiaVoters[userID]
		if exists {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("You already voted for kill\n"),
			}
			return
		}
		if IDPlayer[userID].role != pb.GameAction_MAFIA {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("You're not mafia\n"),
			}
			return
		}
		if IDPlayer[userID].status == DEAD {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("You're dead, you can't kill\n"),
			}
			return
		}
		if IsDay {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("You can kill only during the night\n"),
			}
			return
		}
		victim := action.GetData()
		victimID, exists := PlayerID[victim]
		if !exists {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("Such user does not exists\n"),
			}
			return
		}
		MafiaVoters[userID] = victimID
		MafiaVotes[victimID] = struct{}{}
		if len(MafiaVoters) == MafiaAlive {
			if len(MafiaVotes) == 1 {
				for id := range MafiaVotes {
					killPlayer(id)
					if IDPlayer[id].role == pb.GameAction_MAFIA {
						MafiaAlive--
					} else {
						CitizenAlive--
					}
				}
			} else {
				nobodyKilled()
			}
			IsDay = true
			if MafiaAlive == 0 {
				citizenWin()
			} else if CitizenAlive == MafiaAlive {
				mafiaWin()
			}
			MafiaVoters = map[uuid.UUID]uuid.UUID{}
			MafiaVotes = map[uuid.UUID]struct{}{}
		}
	case pb.GameAction_VOTE:
		_, exists := AllVoters[userID]
		if exists {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("You already voted\n"),
			}
			return
		}
		if IDPlayer[userID].status == DEAD {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("You're dead, you can't vote\n"),
			}
			return
		}
		if !IsDay {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("You can vote only during the day\n"),
			}
			return
		}
		victim := action.GetData()
		victimID, exists := PlayerID[victim]
		if !exists {
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_INFO,
				Data:   stringPtr("Such user does not exists\n"),
			}
			return
		}
		AllVoters[userID] = victimID
		AllVotes[victimID] = AllVotes[victimID] + 1
		if len(AllVoters) == MafiaAlive+CitizenAlive {
			flag := true
			for id, count := range AllVotes {
				if count >= VoteAccept() {
					flag = false
					killPlayer(id)
					if IDPlayer[id].role == pb.GameAction_MAFIA {
						MafiaAlive--
					} else {
						CitizenAlive--
					}
					break
				}
			}
			if flag {
				nobodyKilled()
			}

			IsDay = false
			if MafiaAlive == 0 {
				citizenWin()
			} else if CitizenAlive == MafiaAlive {
				mafiaWin()
			}
			AllVoters = map[uuid.UUID]uuid.UUID{}
			AllVotes = map[uuid.UUID]int{}
		}
	}
}

func nobodyKilled() {
	for _, player := range IDPlayer {
		player.actions <- &pb.GameAction{
			Action: pb.GameAction_KILL,
		}
	}
}

func citizenWin() {
	for _, player := range IDPlayer {
		player.actions <- &pb.GameAction{
			Action: pb.GameAction_CIT_WIN,
		}
	}
	stopGame()
}

func mafiaWin() {
	for _, player := range IDPlayer {
		player.actions <- &pb.GameAction{
			Action: pb.GameAction_MAF_WIN,
		}
	}
	stopGame()
}

func stopGame() {
	GameStarted = false
	InGame = map[uuid.UUID]struct{}{}
	for _, player := range IDPlayer {
		player.status = ALIVE
		if player.players != nil {
			close(player.players)
		}
		if player.chat != nil {
			close(player.chat)
		}
		if player.actions != nil {
			close(player.actions)
		}
	}
}

func killPlayer(userID uuid.UUID) {
	IDPlayer[userID].status = DEAD
	for _, player := range IDPlayer {
		player.actions <- &pb.GameAction{
			Action: pb.GameAction_KILL,
			UserID: convertUUIDtoPB(userID),
		}
	}
}

func sendRoles() {
	var ids []uuid.UUID
	for id := range InGame {
		ids = append(ids, id)
	}
	rand.Seed(time.Now().Unix())
	mafiaIdx := map[int]struct{}{}
	var idx int
	for i := 0; i < MAFIA_COUNT; i++ {
		exists := true
		for exists {
			idx = rand.Intn(MAFIA_COUNT + CITIZEN_COUNT)
			_, exists = mafiaIdx[idx]
		}
		mafiaIdx[idx] = struct{}{}
	}
	for idx, userID := range ids {
		_, exists := mafiaIdx[idx]
		if exists {
			IDPlayer[userID].role = pb.GameAction_MAFIA
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_ROLE,
				Role:   &IDPlayer[userID].role,
			}
			log.Printf("Sended mafia to %s", IDPlayer[userID].username)
		} else {
			IDPlayer[userID].role = pb.GameAction_CITIZEN
			IDPlayer[userID].actions <- &pb.GameAction{
				Action: pb.GameAction_ROLE,
				Role:   &IDPlayer[userID].role,
			}
			log.Printf("Sended citizen to %s", IDPlayer[userID].username)
		}
	}
}

func (s *Server) GetPlayers(id *pb.UUID, server pb.Mafia_GetPlayersServer) error {
	userID := convertPBtoUUID(id)
	players := make(chan *pb.PlayersChange)
	IDPlayer[userID].players = players
	for {
		player, more := <-players
		if !more {
			return nil
		}
		err := server.Send(player)
		if err != nil {
			if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
				log.Printf("Error on getting players change: %v\n", err)
			}
			return err
		}
	}
}

func (s *Server) StartChat(server pb.Mafia_StartChatServer) error {
	msg, err := server.Recv()
	if err != nil {
		if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
			log.Printf("Error on getting chat message: %v", err)
		}
		return nil
	}
	userID := convertPBtoUUID(msg.GetUserID())
	chat := make(chan *pb.ChatMessage)
	IDPlayer[userID].chat = chat
	go func() {
		for {
			msg, more := <-chat
			if !more {
				return
			}
			err := server.Send(msg)
			if err != nil {
				if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
					log.Printf("Error on getting chat message: %v", err)
				}
				return
			}
		}
	}()
	for {
		msg, err := server.Recv()
		if err != nil {
			if !errors.Is(err, io.EOF) && !(status.Convert(err).Message() == "EOF") {
				log.Printf("Error on getting players change: %v", err)
			}
			return nil
		}
		playerStatus := IDPlayer[convertPBtoUUID(msg.GetUserID())].status
		playerRole := IDPlayer[convertPBtoUUID(msg.GetUserID())].role
		for id := range InGame {
			if (playerStatus == DEAD && IDPlayer[id].status == DEAD) ||
				(playerStatus == ALIVE && (IsDay || (playerRole == pb.GameAction_MAFIA && IDPlayer[id].role == pb.GameAction_MAFIA))) {
				IDPlayer[id].chat <- msg
			}
		}
	}
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
