// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: main.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GameAction_ActionType int32

const (
	GameAction_ENTER   GameAction_ActionType = 0
	GameAction_LEAVE   GameAction_ActionType = 1
	GameAction_ROLE    GameAction_ActionType = 2
	GameAction_KILL    GameAction_ActionType = 3
	GameAction_VOTE    GameAction_ActionType = 4
	GameAction_CIT_WIN GameAction_ActionType = 5
	GameAction_MAF_WIN GameAction_ActionType = 6
	GameAction_INFO    GameAction_ActionType = 99
)

// Enum value maps for GameAction_ActionType.
var (
	GameAction_ActionType_name = map[int32]string{
		0:  "ENTER",
		1:  "LEAVE",
		2:  "ROLE",
		3:  "KILL",
		4:  "VOTE",
		5:  "CIT_WIN",
		6:  "MAF_WIN",
		99: "INFO",
	}
	GameAction_ActionType_value = map[string]int32{
		"ENTER":   0,
		"LEAVE":   1,
		"ROLE":    2,
		"KILL":    3,
		"VOTE":    4,
		"CIT_WIN": 5,
		"MAF_WIN": 6,
		"INFO":    99,
	}
)

func (x GameAction_ActionType) Enum() *GameAction_ActionType {
	p := new(GameAction_ActionType)
	*p = x
	return p
}

func (x GameAction_ActionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameAction_ActionType) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[0].Descriptor()
}

func (GameAction_ActionType) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[0]
}

func (x GameAction_ActionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameAction_ActionType.Descriptor instead.
func (GameAction_ActionType) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{5, 0}
}

type GameAction_Roles int32

const (
	GameAction_MAFIA   GameAction_Roles = 0
	GameAction_CITIZEN GameAction_Roles = 1
)

// Enum value maps for GameAction_Roles.
var (
	GameAction_Roles_name = map[int32]string{
		0: "MAFIA",
		1: "CITIZEN",
	}
	GameAction_Roles_value = map[string]int32{
		"MAFIA":   0,
		"CITIZEN": 1,
	}
)

func (x GameAction_Roles) Enum() *GameAction_Roles {
	p := new(GameAction_Roles)
	*p = x
	return p
}

func (x GameAction_Roles) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameAction_Roles) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[1].Descriptor()
}

func (GameAction_Roles) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[1]
}

func (x GameAction_Roles) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameAction_Roles.Descriptor instead.
func (GameAction_Roles) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{5, 1}
}

type PlayersChange_ChangeType int32

const (
	PlayersChange_ADD    PlayersChange_ChangeType = 0
	PlayersChange_REMOVE PlayersChange_ChangeType = 1
)

// Enum value maps for PlayersChange_ChangeType.
var (
	PlayersChange_ChangeType_name = map[int32]string{
		0: "ADD",
		1: "REMOVE",
	}
	PlayersChange_ChangeType_value = map[string]int32{
		"ADD":    0,
		"REMOVE": 1,
	}
)

func (x PlayersChange_ChangeType) Enum() *PlayersChange_ChangeType {
	p := new(PlayersChange_ChangeType)
	*p = x
	return p
}

func (x PlayersChange_ChangeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PlayersChange_ChangeType) Descriptor() protoreflect.EnumDescriptor {
	return file_main_proto_enumTypes[2].Descriptor()
}

func (PlayersChange_ChangeType) Type() protoreflect.EnumType {
	return &file_main_proto_enumTypes[2]
}

func (x PlayersChange_ChangeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PlayersChange_ChangeType.Descriptor instead.
func (PlayersChange_ChangeType) EnumDescriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{6, 0}
}

type NoArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoArgs) Reset() {
	*x = NoArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoArgs) ProtoMessage() {}

func (x *NoArgs) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoArgs.ProtoReflect.Descriptor instead.
func (*NoArgs) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{0}
}

type NoReturn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoReturn) Reset() {
	*x = NoReturn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoReturn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoReturn) ProtoMessage() {}

func (x *NoReturn) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoReturn.ProtoReflect.Descriptor instead.
func (*NoReturn) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{1}
}

type UUID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *UUID) Reset() {
	*x = UUID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UUID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UUID) ProtoMessage() {}

func (x *UUID) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UUID.ProtoReflect.Descriptor instead.
func (*UUID) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{2}
}

func (x *UUID) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type Username struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	UserID   *UUID  `protobuf:"bytes,2,opt,name=UserID,proto3,oneof" json:"UserID,omitempty"`
}

func (x *Username) Reset() {
	*x = Username{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Username) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Username) ProtoMessage() {}

func (x *Username) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Username.ProtoReflect.Descriptor instead.
func (*Username) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{3}
}

func (x *Username) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Username) GetUserID() *UUID {
	if x != nil {
		return x.UserID
	}
	return nil
}

type ListRoomsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rooms []*ListRoomsResponse_Room `protobuf:"bytes,1,rep,name=Rooms,proto3" json:"Rooms,omitempty"`
}

func (x *ListRoomsResponse) Reset() {
	*x = ListRoomsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRoomsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRoomsResponse) ProtoMessage() {}

func (x *ListRoomsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRoomsResponse.ProtoReflect.Descriptor instead.
func (*ListRoomsResponse) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{4}
}

func (x *ListRoomsResponse) GetRooms() []*ListRoomsResponse_Room {
	if x != nil {
		return x.Rooms
	}
	return nil
}

type GameAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action GameAction_ActionType `protobuf:"varint,1,opt,name=Action,proto3,enum=GameAction_ActionType" json:"Action,omitempty"`
	UserID *UUID                 `protobuf:"bytes,2,opt,name=UserID,proto3,oneof" json:"UserID,omitempty"`
	Data   *string               `protobuf:"bytes,3,opt,name=Data,proto3,oneof" json:"Data,omitempty"`
	Role   *GameAction_Roles     `protobuf:"varint,4,opt,name=Role,proto3,enum=GameAction_Roles,oneof" json:"Role,omitempty"`
}

func (x *GameAction) Reset() {
	*x = GameAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameAction) ProtoMessage() {}

func (x *GameAction) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameAction.ProtoReflect.Descriptor instead.
func (*GameAction) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{5}
}

func (x *GameAction) GetAction() GameAction_ActionType {
	if x != nil {
		return x.Action
	}
	return GameAction_ENTER
}

func (x *GameAction) GetUserID() *UUID {
	if x != nil {
		return x.UserID
	}
	return nil
}

func (x *GameAction) GetData() string {
	if x != nil && x.Data != nil {
		return *x.Data
	}
	return ""
}

func (x *GameAction) GetRole() GameAction_Roles {
	if x != nil && x.Role != nil {
		return *x.Role
	}
	return GameAction_MAFIA
}

type PlayersChange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action   PlayersChange_ChangeType `protobuf:"varint,1,opt,name=Action,proto3,enum=PlayersChange_ChangeType" json:"Action,omitempty"`
	UserID   *UUID                    `protobuf:"bytes,2,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Username *string                  `protobuf:"bytes,3,opt,name=Username,proto3,oneof" json:"Username,omitempty"`
}

func (x *PlayersChange) Reset() {
	*x = PlayersChange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayersChange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayersChange) ProtoMessage() {}

func (x *PlayersChange) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayersChange.ProtoReflect.Descriptor instead.
func (*PlayersChange) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{6}
}

func (x *PlayersChange) GetAction() PlayersChange_ChangeType {
	if x != nil {
		return x.Action
	}
	return PlayersChange_ADD
}

func (x *PlayersChange) GetUserID() *UUID {
	if x != nil {
		return x.UserID
	}
	return nil
}

func (x *PlayersChange) GetUsername() string {
	if x != nil && x.Username != nil {
		return *x.Username
	}
	return ""
}

type ChatMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID *UUID  `protobuf:"bytes,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Text   string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{7}
}

func (x *ChatMessage) GetUserID() *UUID {
	if x != nil {
		return x.UserID
	}
	return nil
}

func (x *ChatMessage) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type ListRoomsResponse_Room struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	UserNum uint32 `protobuf:"varint,2,opt,name=UserNum,proto3" json:"UserNum,omitempty"`
}

func (x *ListRoomsResponse_Room) Reset() {
	*x = ListRoomsResponse_Room{}
	if protoimpl.UnsafeEnabled {
		mi := &file_main_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRoomsResponse_Room) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRoomsResponse_Room) ProtoMessage() {}

func (x *ListRoomsResponse_Room) ProtoReflect() protoreflect.Message {
	mi := &file_main_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRoomsResponse_Room.ProtoReflect.Descriptor instead.
func (*ListRoomsResponse_Room) Descriptor() ([]byte, []int) {
	return file_main_proto_rawDescGZIP(), []int{4, 0}
}

func (x *ListRoomsResponse_Room) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListRoomsResponse_Room) GetUserNum() uint32 {
	if x != nil {
		return x.UserNum
	}
	return 0
}

var File_main_proto protoreflect.FileDescriptor

var file_main_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x08, 0x0a, 0x06,
	0x4e, 0x6f, 0x41, 0x72, 0x67, 0x73, 0x22, 0x0a, 0x0a, 0x08, 0x4e, 0x6f, 0x52, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x22, 0x1c, 0x0a, 0x04, 0x55, 0x55, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x55, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x48,
	0x00, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07,
	0x5f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x78, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x6f, 0x6f, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x05,
	0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x05, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x1a, 0x34, 0x0a, 0x04, 0x52,
	0x6f, 0x6f, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x4e,
	0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x75,
	0x6d, 0x22, 0xc9, 0x02, 0x0a, 0x0a, 0x47, 0x61, 0x6d, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x2e, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x16, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x22, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x48, 0x00, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x01, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a,
	0x04, 0x52, 0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x47, 0x61,
	0x6d, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x48, 0x02,
	0x52, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x88, 0x01, 0x01, 0x22, 0x64, 0x0a, 0x0a, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x4e, 0x54, 0x45, 0x52,
	0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4c, 0x45, 0x41, 0x56, 0x45, 0x10, 0x01, 0x12, 0x08, 0x0a,
	0x04, 0x52, 0x4f, 0x4c, 0x45, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x4b, 0x49, 0x4c, 0x4c, 0x10,
	0x03, 0x12, 0x08, 0x0a, 0x04, 0x56, 0x4f, 0x54, 0x45, 0x10, 0x04, 0x12, 0x0b, 0x0a, 0x07, 0x43,
	0x49, 0x54, 0x5f, 0x57, 0x49, 0x4e, 0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x41, 0x46, 0x5f,
	0x57, 0x49, 0x4e, 0x10, 0x06, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x63, 0x22,
	0x1f, 0x0a, 0x05, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x41, 0x46, 0x49,
	0x41, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x49, 0x54, 0x49, 0x5a, 0x45, 0x4e, 0x10, 0x01,
	0x42, 0x09, 0x0a, 0x07, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x44, 0x61, 0x74, 0x61, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0xb2, 0x01,
	0x0a, 0x0d, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12,
	0x31, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x19, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2e,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x44, 0x12, 0x1f, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x88,
	0x01, 0x01, 0x22, 0x21, 0x0a, 0x0a, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x07, 0x0a, 0x03, 0x41, 0x44, 0x44, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x45, 0x4d,
	0x4f, 0x56, 0x45, 0x10, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x22, 0x40, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x1d, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x32, 0xda, 0x01, 0x0a, 0x05, 0x4d, 0x61, 0x66, 0x69, 0x61, 0x12, 0x21,
	0x0a, 0x0b, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x09, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x1a, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x22,
	0x00, 0x12, 0x2a, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x12, 0x07,
	0x2e, 0x4e, 0x6f, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x12, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x6f,
	0x6f, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2a, 0x0a,
	0x08, 0x50, 0x6c, 0x61, 0x79, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x0b, 0x2e, 0x47, 0x61, 0x6d, 0x65,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x0b, 0x2e, 0x47, 0x61, 0x6d, 0x65, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x12, 0x27, 0x0a, 0x0a, 0x47, 0x65, 0x74,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x12, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x1a, 0x0e,
	0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x22, 0x00,
	0x30, 0x01, 0x12, 0x2d, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x68, 0x61, 0x74, 0x12,
	0x0c, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x0c, 0x2e,
	0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30,
	0x01, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_main_proto_rawDescOnce sync.Once
	file_main_proto_rawDescData = file_main_proto_rawDesc
)

func file_main_proto_rawDescGZIP() []byte {
	file_main_proto_rawDescOnce.Do(func() {
		file_main_proto_rawDescData = protoimpl.X.CompressGZIP(file_main_proto_rawDescData)
	})
	return file_main_proto_rawDescData
}

var file_main_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_main_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_main_proto_goTypes = []interface{}{
	(GameAction_ActionType)(0),     // 0: GameAction.ActionType
	(GameAction_Roles)(0),          // 1: GameAction.Roles
	(PlayersChange_ChangeType)(0),  // 2: PlayersChange.ChangeType
	(*NoArgs)(nil),                 // 3: NoArgs
	(*NoReturn)(nil),               // 4: NoReturn
	(*UUID)(nil),                   // 5: UUID
	(*Username)(nil),               // 6: Username
	(*ListRoomsResponse)(nil),      // 7: ListRoomsResponse
	(*GameAction)(nil),             // 8: GameAction
	(*PlayersChange)(nil),          // 9: PlayersChange
	(*ChatMessage)(nil),            // 10: ChatMessage
	(*ListRoomsResponse_Room)(nil), // 11: ListRoomsResponse.Room
}
var file_main_proto_depIdxs = []int32{
	5,  // 0: Username.UserID:type_name -> UUID
	11, // 1: ListRoomsResponse.Rooms:type_name -> ListRoomsResponse.Room
	0,  // 2: GameAction.Action:type_name -> GameAction.ActionType
	5,  // 3: GameAction.UserID:type_name -> UUID
	1,  // 4: GameAction.Role:type_name -> GameAction.Roles
	2,  // 5: PlayersChange.Action:type_name -> PlayersChange.ChangeType
	5,  // 6: PlayersChange.UserID:type_name -> UUID
	5,  // 7: ChatMessage.UserID:type_name -> UUID
	6,  // 8: Mafia.SetUsername:input_type -> Username
	3,  // 9: Mafia.ListRooms:input_type -> NoArgs
	8,  // 10: Mafia.PlayGame:input_type -> GameAction
	5,  // 11: Mafia.GetPlayers:input_type -> UUID
	10, // 12: Mafia.StartChat:input_type -> ChatMessage
	5,  // 13: Mafia.SetUsername:output_type -> UUID
	7,  // 14: Mafia.ListRooms:output_type -> ListRoomsResponse
	8,  // 15: Mafia.PlayGame:output_type -> GameAction
	9,  // 16: Mafia.GetPlayers:output_type -> PlayersChange
	10, // 17: Mafia.StartChat:output_type -> ChatMessage
	13, // [13:18] is the sub-list for method output_type
	8,  // [8:13] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_main_proto_init() }
func file_main_proto_init() {
	if File_main_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_main_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoReturn); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UUID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Username); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRoomsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameAction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayersChange); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_main_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRoomsResponse_Room); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_main_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_main_proto_msgTypes[5].OneofWrappers = []interface{}{}
	file_main_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_main_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_main_proto_goTypes,
		DependencyIndexes: file_main_proto_depIdxs,
		EnumInfos:         file_main_proto_enumTypes,
		MessageInfos:      file_main_proto_msgTypes,
	}.Build()
	File_main_proto = out.File
	file_main_proto_rawDesc = nil
	file_main_proto_goTypes = nil
	file_main_proto_depIdxs = nil
}