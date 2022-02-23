package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
)

type MsgType uint8

const (
	PORT = 8080
)

const (
	CONNECT MsgType = iota
	NEW_USER
	DISCONNECT
	ROOM_LIST
	ROOM_CONNECT
	ROOM_DISCONNECT
	TEXT
	ERROR
)

type Message struct {
	Type   MsgType
	UserID uuid.NullUUID
	Data   []byte
}

func (msg *Message) Encode() ([]byte, error) {
	var payload bytes.Buffer
	var uuidRaw []byte
	if msg.UserID.Valid {
		uuidRaw, _ = msg.UserID.MarshalBinary()
	} else {
		uuidRaw = make([]byte, 16)
	}
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(1+len(uuidRaw)+len(msg.Data)))
	payload.Write(size)
	payload.WriteByte(byte(msg.Type))
	payload.Write(uuidRaw)
	payload.Write(msg.Data)
	return payload.Bytes(), nil
}

func (msg *Message) Decode(data []byte) error {
	payload := bytes.NewBuffer(data)
	b, _ := payload.ReadByte()
	msg.Type = MsgType(b)
	uuidRaw := make([]byte, 16)
	payload.Read(uuidRaw)
	msg.UserID.UnmarshalBinary(uuidRaw)
	msg.Data = make([]byte, payload.Len())
	payload.Read(msg.Data)
	return nil
}
