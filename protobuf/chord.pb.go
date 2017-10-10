// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chord.proto

/*
Package chordMsgs is a generated protocol buffer package.

It is generated from these files:
	chord.proto

It has these top-level messages:
	SendIdMessage
	FingerMessage
	PredMessage
	SendFingersMessage
	ChordMessage
	NetworkMessage
*/
package chordMsgs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ChordMessage_Command int32

const (
	ChordMessage_Ping       ChordMessage_Command = 1
	ChordMessage_Pong       ChordMessage_Command = 2
	ChordMessage_GetPred    ChordMessage_Command = 3
	ChordMessage_GetId      ChordMessage_Command = 4
	ChordMessage_GetFingers ChordMessage_Command = 5
	ChordMessage_ClaimPred  ChordMessage_Command = 6
	ChordMessage_GetSucc    ChordMessage_Command = 7
)

var ChordMessage_Command_name = map[int32]string{
	1: "Ping",
	2: "Pong",
	3: "GetPred",
	4: "GetId",
	5: "GetFingers",
	6: "ClaimPred",
	7: "GetSucc",
}
var ChordMessage_Command_value = map[string]int32{
	"Ping":       1,
	"Pong":       2,
	"GetPred":    3,
	"GetId":      4,
	"GetFingers": 5,
	"ClaimPred":  6,
	"GetSucc":    7,
}

func (x ChordMessage_Command) Enum() *ChordMessage_Command {
	p := new(ChordMessage_Command)
	*p = x
	return p
}
func (x ChordMessage_Command) String() string {
	return proto.EnumName(ChordMessage_Command_name, int32(x))
}
func (x *ChordMessage_Command) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ChordMessage_Command_value, data, "ChordMessage_Command")
	if err != nil {
		return err
	}
	*x = ChordMessage_Command(value)
	return nil
}
func (ChordMessage_Command) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{4, 0} }

type SendIdMessage struct {
	Id               *string `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SendIdMessage) Reset()                    { *m = SendIdMessage{} }
func (m *SendIdMessage) String() string            { return proto.CompactTextString(m) }
func (*SendIdMessage) ProtoMessage()               {}
func (*SendIdMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SendIdMessage) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

type FingerMessage struct {
	Address          *string `protobuf:"bytes,1,req,name=address" json:"address,omitempty"`
	Id               *string `protobuf:"bytes,2,req,name=id" json:"id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *FingerMessage) Reset()                    { *m = FingerMessage{} }
func (m *FingerMessage) String() string            { return proto.CompactTextString(m) }
func (*FingerMessage) ProtoMessage()               {}
func (*FingerMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *FingerMessage) GetAddress() string {
	if m != nil && m.Address != nil {
		return *m.Address
	}
	return ""
}

func (m *FingerMessage) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

type PredMessage struct {
	Pred             *FingerMessage `protobuf:"bytes,1,req,name=pred" json:"pred,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *PredMessage) Reset()                    { *m = PredMessage{} }
func (m *PredMessage) String() string            { return proto.CompactTextString(m) }
func (*PredMessage) ProtoMessage()               {}
func (*PredMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PredMessage) GetPred() *FingerMessage {
	if m != nil {
		return m.Pred
	}
	return nil
}

type SendFingersMessage struct {
	Fingers          []*FingerMessage `protobuf:"bytes,1,rep,name=fingers" json:"fingers,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *SendFingersMessage) Reset()                    { *m = SendFingersMessage{} }
func (m *SendFingersMessage) String() string            { return proto.CompactTextString(m) }
func (*SendFingersMessage) ProtoMessage()               {}
func (*SendFingersMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SendFingersMessage) GetFingers() []*FingerMessage {
	if m != nil {
		return m.Fingers
	}
	return nil
}

type ChordMessage struct {
	Cmd              *ChordMessage_Command `protobuf:"varint,1,req,name=cmd,enum=chordMsgs.ChordMessage_Command" json:"cmd,omitempty"`
	Cpmsg            *PredMessage          `protobuf:"bytes,2,opt,name=cpmsg" json:"cpmsg,omitempty"`
	Sidmsg           *SendIdMessage        `protobuf:"bytes,4,opt,name=sidmsg" json:"sidmsg,omitempty"`
	Sfmsg            *SendFingersMessage   `protobuf:"bytes,5,opt,name=sfmsg" json:"sfmsg,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (m *ChordMessage) Reset()                    { *m = ChordMessage{} }
func (m *ChordMessage) String() string            { return proto.CompactTextString(m) }
func (*ChordMessage) ProtoMessage()               {}
func (*ChordMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ChordMessage) GetCmd() ChordMessage_Command {
	if m != nil && m.Cmd != nil {
		return *m.Cmd
	}
	return ChordMessage_Ping
}

func (m *ChordMessage) GetCpmsg() *PredMessage {
	if m != nil {
		return m.Cpmsg
	}
	return nil
}

func (m *ChordMessage) GetSidmsg() *SendIdMessage {
	if m != nil {
		return m.Sidmsg
	}
	return nil
}

func (m *ChordMessage) GetSfmsg() *SendFingersMessage {
	if m != nil {
		return m.Sfmsg
	}
	return nil
}

type NetworkMessage struct {
	Proto            *uint32 `protobuf:"varint,1,opt,name=proto" json:"proto,omitempty"`
	Msg              *string `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *NetworkMessage) Reset()                    { *m = NetworkMessage{} }
func (m *NetworkMessage) String() string            { return proto.CompactTextString(m) }
func (*NetworkMessage) ProtoMessage()               {}
func (*NetworkMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *NetworkMessage) GetProto() uint32 {
	if m != nil && m.Proto != nil {
		return *m.Proto
	}
	return 0
}

func (m *NetworkMessage) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*SendIdMessage)(nil), "chordMsgs.SendIdMessage")
	proto.RegisterType((*FingerMessage)(nil), "chordMsgs.FingerMessage")
	proto.RegisterType((*PredMessage)(nil), "chordMsgs.PredMessage")
	proto.RegisterType((*SendFingersMessage)(nil), "chordMsgs.SendFingersMessage")
	proto.RegisterType((*ChordMessage)(nil), "chordMsgs.ChordMessage")
	proto.RegisterType((*NetworkMessage)(nil), "chordMsgs.NetworkMessage")
	proto.RegisterEnum("chordMsgs.ChordMessage_Command", ChordMessage_Command_name, ChordMessage_Command_value)
}

func init() { proto.RegisterFile("chord.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 346 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x41, 0x6b, 0xbb, 0x30,
	0x18, 0xc6, 0x31, 0x6a, 0xfd, 0xfb, 0xfa, 0x57, 0x42, 0x18, 0xc3, 0xcb, 0xa8, 0x78, 0xea, 0xa1,
	0xc8, 0xe6, 0x2e, 0x1b, 0x3b, 0x16, 0xd6, 0xf5, 0xd0, 0x51, 0xd2, 0x4f, 0x20, 0x26, 0x75, 0xb2,
	0xa9, 0x25, 0x71, 0xec, 0xb3, 0xed, 0xdb, 0x0d, 0x93, 0x58, 0x74, 0x87, 0xdd, 0xde, 0x57, 0x7f,
	0xcf, 0x93, 0x27, 0x0f, 0x81, 0xa0, 0x7c, 0xeb, 0x04, 0xcb, 0xce, 0xa2, 0xeb, 0x3b, 0xe2, 0xab,
	0x65, 0x2f, 0x2b, 0x99, 0x2e, 0x21, 0x3c, 0xf2, 0x96, 0xed, 0xd8, 0x9e, 0x4b, 0x59, 0x54, 0x9c,
	0x44, 0x80, 0x6a, 0x16, 0x5b, 0x09, 0x5a, 0xf9, 0x14, 0xd5, 0x2c, 0x7d, 0x84, 0xf0, 0xb9, 0x6e,
	0x2b, 0x2e, 0x46, 0x20, 0x06, 0xaf, 0x60, 0x4c, 0x70, 0x29, 0x0d, 0x35, 0xae, 0x46, 0x8a, 0x2e,
	0xd2, 0x27, 0x08, 0x0e, 0x82, 0x5f, 0x9c, 0xd7, 0xe0, 0x9c, 0x05, 0xd7, 0xde, 0x41, 0x1e, 0x67,
	0x97, 0x10, 0xd9, 0xec, 0x00, 0xaa, 0xa8, 0xf4, 0x05, 0xc8, 0x10, 0x4c, 0xff, 0x92, 0xa3, 0x47,
	0x0e, 0xde, 0x49, 0x7f, 0x89, 0xad, 0xc4, 0xfe, 0xd3, 0x66, 0x04, 0xd3, 0x6f, 0x04, 0xff, 0x37,
	0x0a, 0x32, 0x26, 0x77, 0x60, 0x97, 0x8d, 0xce, 0x11, 0xe5, 0xcb, 0x89, 0xc1, 0x94, 0xca, 0x36,
	0x5d, 0xd3, 0x14, 0x2d, 0xa3, 0x03, 0x4b, 0xd6, 0xe0, 0x96, 0xe7, 0x46, 0x56, 0x31, 0x4a, 0xac,
	0x55, 0x90, 0x5f, 0x4f, 0x44, 0x93, 0x2b, 0x52, 0x0d, 0x91, 0x5b, 0x58, 0xc8, 0x9a, 0x0d, 0xb8,
	0xa3, 0xf0, 0x69, 0xc8, 0x59, 0xdb, 0xd4, 0x70, 0xe4, 0x1e, 0x5c, 0x79, 0x1a, 0x04, 0xae, 0x12,
	0xdc, 0xfc, 0x12, 0xcc, 0x5b, 0xa0, 0x9a, 0x4d, 0x0b, 0xf0, 0x4c, 0x48, 0xf2, 0x0f, 0x9c, 0x43,
	0xdd, 0x56, 0xd8, 0x52, 0x53, 0xd7, 0x56, 0x18, 0x91, 0x00, 0xbc, 0x2d, 0xef, 0x87, 0x78, 0xd8,
	0x26, 0x3e, 0xb8, 0x5b, 0xde, 0xef, 0x18, 0x76, 0x48, 0x04, 0xb0, 0xe5, 0xbd, 0xb1, 0xc4, 0x2e,
	0x09, 0xc1, 0xdf, 0x7c, 0x14, 0x75, 0xa3, 0xc8, 0x85, 0x91, 0x1d, 0x3f, 0xcb, 0x12, 0x7b, 0xe9,
	0x03, 0x44, 0xaf, 0xbc, 0xff, 0xea, 0xc4, 0xfb, 0x58, 0xde, 0x15, 0xb8, 0xea, 0x11, 0xc5, 0x56,
	0x62, 0xad, 0x42, 0xaa, 0x17, 0x82, 0xc1, 0x1e, 0xdb, 0xf1, 0xe9, 0x30, 0xfe, 0x04, 0x00, 0x00,
	0xff, 0xff, 0x60, 0x9a, 0xa6, 0x5d, 0x71, 0x02, 0x00, 0x00,
}
