// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peer/events.proto

package peer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "github.com/hyperledger/fabric/protos/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type EventType int32

const (
	EventType_REGISTER  EventType = 0
	EventType_BLOCK     EventType = 1
	EventType_CHAINCODE EventType = 2
	EventType_REJECTION EventType = 3
)

var EventType_name = map[int32]string{
	0: "REGISTER",
	1: "BLOCK",
	2: "CHAINCODE",
	3: "REJECTION",
}
var EventType_value = map[string]int32{
	"REGISTER":  0,
	"BLOCK":     1,
	"CHAINCODE": 2,
	"REJECTION": 3,
}

func (x EventType) String() string {
	return proto.EnumName(EventType_name, int32(x))
}
func (EventType) EnumDescriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

// ChaincodeReg is used for registering chaincode Interests
// when EventType is CHAINCODE
type ChaincodeReg struct {
	ChaincodeId string `protobuf:"bytes,1,opt,name=chaincode_id,json=chaincodeId" json:"chaincode_id,omitempty"`
	EventName   string `protobuf:"bytes,2,opt,name=event_name,json=eventName" json:"event_name,omitempty"`
}

func (m *ChaincodeReg) Reset()                    { *m = ChaincodeReg{} }
func (m *ChaincodeReg) String() string            { return proto.CompactTextString(m) }
func (*ChaincodeReg) ProtoMessage()               {}
func (*ChaincodeReg) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

func (m *ChaincodeReg) GetChaincodeId() string {
	if m != nil {
		return m.ChaincodeId
	}
	return ""
}

func (m *ChaincodeReg) GetEventName() string {
	if m != nil {
		return m.EventName
	}
	return ""
}

type Interest struct {
	EventType EventType `protobuf:"varint,1,opt,name=event_type,json=eventType,enum=protos.EventType" json:"event_type,omitempty"`
	// Ideally we should just have the following oneof for different
	// Reg types and get rid of EventType. But this is an API change
	// Additional Reg types may add messages specific to their type
	// to the oneof.
	//
	// Types that are valid to be assigned to RegInfo:
	//	*Interest_ChaincodeRegInfo
	RegInfo isInterest_RegInfo `protobuf_oneof:"RegInfo"`
	ChainID string             `protobuf:"bytes,3,opt,name=chainID" json:"chainID,omitempty"`
}

func (m *Interest) Reset()                    { *m = Interest{} }
func (m *Interest) String() string            { return proto.CompactTextString(m) }
func (*Interest) ProtoMessage()               {}
func (*Interest) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{1} }

type isInterest_RegInfo interface {
	isInterest_RegInfo()
}

type Interest_ChaincodeRegInfo struct {
	ChaincodeRegInfo *ChaincodeReg `protobuf:"bytes,2,opt,name=chaincode_reg_info,json=chaincodeRegInfo,oneof"`
}

func (*Interest_ChaincodeRegInfo) isInterest_RegInfo() {}

func (m *Interest) GetRegInfo() isInterest_RegInfo {
	if m != nil {
		return m.RegInfo
	}
	return nil
}

func (m *Interest) GetEventType() EventType {
	if m != nil {
		return m.EventType
	}
	return EventType_REGISTER
}

func (m *Interest) GetChaincodeRegInfo() *ChaincodeReg {
	if x, ok := m.GetRegInfo().(*Interest_ChaincodeRegInfo); ok {
		return x.ChaincodeRegInfo
	}
	return nil
}

func (m *Interest) GetChainID() string {
	if m != nil {
		return m.ChainID
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Interest) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Interest_OneofMarshaler, _Interest_OneofUnmarshaler, _Interest_OneofSizer, []interface{}{
		(*Interest_ChaincodeRegInfo)(nil),
	}
}

func _Interest_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Interest)
	// RegInfo
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeRegInfo); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Interest.RegInfo has unexpected type %T", x)
	}
	return nil
}

func _Interest_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Interest)
	switch tag {
	case 2: // RegInfo.chaincode_reg_info
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeReg)
		err := b.DecodeMessage(msg)
		m.RegInfo = &Interest_ChaincodeRegInfo{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Interest_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Interest)
	// RegInfo
	switch x := m.RegInfo.(type) {
	case *Interest_ChaincodeRegInfo:
		s := proto.Size(x.ChaincodeRegInfo)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// ---------- consumer events ---------
// Register is sent by consumers for registering events
// string type - "register"
type Register struct {
	Events []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *Register) Reset()                    { *m = Register{} }
func (m *Register) String() string            { return proto.CompactTextString(m) }
func (*Register) ProtoMessage()               {}
func (*Register) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{2} }

func (m *Register) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}

// Rejection is sent by consumers for erroneous transaction rejection events
// string type - "rejection"
type Rejection struct {
	Tx       *Transaction `protobuf:"bytes,1,opt,name=tx" json:"tx,omitempty"`
	ErrorMsg string       `protobuf:"bytes,2,opt,name=error_msg,json=errorMsg" json:"error_msg,omitempty"`
}

func (m *Rejection) Reset()                    { *m = Rejection{} }
func (m *Rejection) String() string            { return proto.CompactTextString(m) }
func (*Rejection) ProtoMessage()               {}
func (*Rejection) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{3} }

func (m *Rejection) GetTx() *Transaction {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *Rejection) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}

// ---------- producer events ---------
type Unregister struct {
	Events []*Interest `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *Unregister) Reset()                    { *m = Unregister{} }
func (m *Unregister) String() string            { return proto.CompactTextString(m) }
func (*Unregister) ProtoMessage()               {}
func (*Unregister) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{4} }

func (m *Unregister) GetEvents() []*Interest {
	if m != nil {
		return m.Events
	}
	return nil
}

// SignedEvent is used for any communication between consumer and producer
type SignedEvent struct {
	// Signature over the event bytes
	Signature []byte `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	// Marshal of Event object
	EventBytes []byte `protobuf:"bytes,2,opt,name=eventBytes,proto3" json:"eventBytes,omitempty"`
}

func (m *SignedEvent) Reset()                    { *m = SignedEvent{} }
func (m *SignedEvent) String() string            { return proto.CompactTextString(m) }
func (*SignedEvent) ProtoMessage()               {}
func (*SignedEvent) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{5} }

func (m *SignedEvent) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *SignedEvent) GetEventBytes() []byte {
	if m != nil {
		return m.EventBytes
	}
	return nil
}

// Event is used by
//  - consumers (adapters) to send Register
//  - producer to advertise supported types and events
type Event struct {
	// Types that are valid to be assigned to Event:
	//	*Event_Register
	//	*Event_Block
	//	*Event_ChaincodeEvent
	//	*Event_Rejection
	//	*Event_Unregister
	Event isEvent_Event `protobuf_oneof:"Event"`
	// Creator of the event, specified as a certificate chain
	Creator []byte `protobuf:"bytes,6,opt,name=creator,proto3" json:"creator,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{6} }

type isEvent_Event interface {
	isEvent_Event()
}

type Event_Register struct {
	Register *Register `protobuf:"bytes,1,opt,name=register,oneof"`
}
type Event_Block struct {
	Block *common.Block `protobuf:"bytes,2,opt,name=block,oneof"`
}
type Event_ChaincodeEvent struct {
	ChaincodeEvent *ChaincodeEvent `protobuf:"bytes,3,opt,name=chaincode_event,json=chaincodeEvent,oneof"`
}
type Event_Rejection struct {
	Rejection *Rejection `protobuf:"bytes,4,opt,name=rejection,oneof"`
}
type Event_Unregister struct {
	Unregister *Unregister `protobuf:"bytes,5,opt,name=unregister,oneof"`
}

func (*Event_Register) isEvent_Event()       {}
func (*Event_Block) isEvent_Event()          {}
func (*Event_ChaincodeEvent) isEvent_Event() {}
func (*Event_Rejection) isEvent_Event()      {}
func (*Event_Unregister) isEvent_Event()     {}

func (m *Event) GetEvent() isEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *Event) GetRegister() *Register {
	if x, ok := m.GetEvent().(*Event_Register); ok {
		return x.Register
	}
	return nil
}

func (m *Event) GetBlock() *common.Block {
	if x, ok := m.GetEvent().(*Event_Block); ok {
		return x.Block
	}
	return nil
}

func (m *Event) GetChaincodeEvent() *ChaincodeEvent {
	if x, ok := m.GetEvent().(*Event_ChaincodeEvent); ok {
		return x.ChaincodeEvent
	}
	return nil
}

func (m *Event) GetRejection() *Rejection {
	if x, ok := m.GetEvent().(*Event_Rejection); ok {
		return x.Rejection
	}
	return nil
}

func (m *Event) GetUnregister() *Unregister {
	if x, ok := m.GetEvent().(*Event_Unregister); ok {
		return x.Unregister
	}
	return nil
}

func (m *Event) GetCreator() []byte {
	if m != nil {
		return m.Creator
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_Register)(nil),
		(*Event_Block)(nil),
		(*Event_ChaincodeEvent)(nil),
		(*Event_Rejection)(nil),
		(*Event_Unregister)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// Event
	switch x := m.Event.(type) {
	case *Event_Register:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Register); err != nil {
			return err
		}
	case *Event_Block:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Block); err != nil {
			return err
		}
	case *Event_ChaincodeEvent:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChaincodeEvent); err != nil {
			return err
		}
	case *Event_Rejection:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Rejection); err != nil {
			return err
		}
	case *Event_Unregister:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Unregister); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.Event has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 1: // Event.register
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Register)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Register{msg}
		return true, err
	case 2: // Event.block
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(common.Block)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Block{msg}
		return true, err
	case 3: // Event.chaincode_event
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChaincodeEvent)
		err := b.DecodeMessage(msg)
		m.Event = &Event_ChaincodeEvent{msg}
		return true, err
	case 4: // Event.rejection
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Rejection)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Rejection{msg}
		return true, err
	case 5: // Event.unregister
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Unregister)
		err := b.DecodeMessage(msg)
		m.Event = &Event_Unregister{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// Event
	switch x := m.Event.(type) {
	case *Event_Register:
		s := proto.Size(x.Register)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Block:
		s := proto.Size(x.Block)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_ChaincodeEvent:
		s := proto.Size(x.ChaincodeEvent)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Rejection:
		s := proto.Size(x.Rejection)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Unregister:
		s := proto.Size(x.Unregister)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ChaincodeReg)(nil), "protos.ChaincodeReg")
	proto.RegisterType((*Interest)(nil), "protos.Interest")
	proto.RegisterType((*Register)(nil), "protos.Register")
	proto.RegisterType((*Rejection)(nil), "protos.Rejection")
	proto.RegisterType((*Unregister)(nil), "protos.Unregister")
	proto.RegisterType((*SignedEvent)(nil), "protos.SignedEvent")
	proto.RegisterType((*Event)(nil), "protos.Event")
	proto.RegisterEnum("protos.EventType", EventType_name, EventType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Events service

type EventsClient interface {
	// event chatting using Event
	Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error)
}

type eventsClient struct {
	cc *grpc.ClientConn
}

func NewEventsClient(cc *grpc.ClientConn) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) Chat(ctx context.Context, opts ...grpc.CallOption) (Events_ChatClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Events_serviceDesc.Streams[0], c.cc, "/protos.Events/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsChatClient{stream}
	return x, nil
}

type Events_ChatClient interface {
	Send(*SignedEvent) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventsChatClient struct {
	grpc.ClientStream
}

func (x *eventsChatClient) Send(m *SignedEvent) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventsChatClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Events service

type EventsServer interface {
	// event chatting using Event
	Chat(Events_ChatServer) error
}

func RegisterEventsServer(s *grpc.Server, srv EventsServer) {
	s.RegisterService(&_Events_serviceDesc, srv)
}

func _Events_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventsServer).Chat(&eventsChatServer{stream})
}

type Events_ChatServer interface {
	Send(*Event) error
	Recv() (*SignedEvent, error)
	grpc.ServerStream
}

type eventsChatServer struct {
	grpc.ServerStream
}

func (x *eventsChatServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventsChatServer) Recv() (*SignedEvent, error) {
	m := new(SignedEvent)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Events_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Events",
	HandlerType: (*EventsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Chat",
			Handler:       _Events_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peer/events.proto",
}

func init() { proto.RegisterFile("peer/events.proto", fileDescriptor5) }

var fileDescriptor5 = []byte{
	// 611 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4f, 0x6f, 0x9b, 0x4e,
	0x10, 0x05, 0x27, 0x76, 0xcc, 0xd8, 0xce, 0xcf, 0xd9, 0xfc, 0x14, 0x21, 0xf7, 0x8f, 0x52, 0xaa,
	0x4a, 0x6e, 0x0f, 0x26, 0xa5, 0x51, 0x0f, 0xbd, 0x05, 0x82, 0x0a, 0x4d, 0xf3, 0x47, 0x1b, 0xf7,
	0xd2, 0x43, 0x2d, 0x8c, 0x27, 0x98, 0x26, 0x06, 0x6b, 0xd9, 0x54, 0xf1, 0x27, 0xea, 0xa9, 0xdf,
	0xb1, 0x62, 0x61, 0x81, 0xb4, 0xa7, 0x9e, 0xf0, 0xcc, 0xbc, 0x37, 0xfb, 0xf6, 0xcd, 0xac, 0x61,
	0x6f, 0x8d, 0xc8, 0x4c, 0xfc, 0x81, 0x09, 0xcf, 0x26, 0x6b, 0x96, 0xf2, 0x94, 0x74, 0xc4, 0x27,
	0x1b, 0xed, 0x87, 0xe9, 0x6a, 0x95, 0x26, 0x66, 0xf1, 0x29, 0x8a, 0xa3, 0x91, 0xc0, 0x87, 0xcb,
	0x20, 0x4e, 0xc2, 0x74, 0x81, 0x33, 0xc1, 0x2c, 0x6b, 0x07, 0xa2, 0xc6, 0x59, 0x90, 0x64, 0x41,
	0xc8, 0x63, 0xc9, 0x31, 0xae, 0xa0, 0xef, 0x48, 0x02, 0xc5, 0x88, 0xbc, 0x80, 0x7e, 0xdd, 0x20,
	0x5e, 0xe8, 0xea, 0xa1, 0x3a, 0xd6, 0x68, 0xaf, 0xca, 0xf9, 0x0b, 0xf2, 0x0c, 0x40, 0x74, 0x9e,
	0x25, 0xc1, 0x0a, 0xf5, 0x96, 0x00, 0x68, 0x22, 0x73, 0x11, 0xac, 0xd0, 0xf8, 0xa9, 0x42, 0xd7,
	0x4f, 0x38, 0x32, 0xcc, 0x38, 0x39, 0x92, 0x58, 0xbe, 0x59, 0xa3, 0x68, 0xb6, 0x6b, 0xed, 0x15,
	0x47, 0x67, 0x13, 0x37, 0xaf, 0x4c, 0x37, 0x6b, 0x2c, 0xe9, 0xf9, 0x4f, 0x72, 0x0a, 0xa4, 0x16,
	0xc0, 0x30, 0x9a, 0xc5, 0xc9, 0x4d, 0x2a, 0x4e, 0xe9, 0x59, 0xff, 0x4b, 0x66, 0x53, 0xb2, 0xa7,
	0xd0, 0x61, 0xd8, 0x88, 0xfd, 0xe4, 0x26, 0x25, 0x3a, 0xec, 0x88, 0x9c, 0x7f, 0xaa, 0x6f, 0x09,
	0x81, 0x32, 0xb4, 0x35, 0xd8, 0x29, 0x41, 0xc6, 0x31, 0x74, 0x29, 0x46, 0x71, 0xc6, 0x91, 0x91,
	0x31, 0x74, 0x0a, 0xa3, 0x75, 0xf5, 0x70, 0x6b, 0xdc, 0xb3, 0x86, 0xf2, 0x28, 0x79, 0x15, 0x5a,
	0xd6, 0x8d, 0x73, 0xd0, 0x28, 0x7e, 0x47, 0x61, 0x22, 0x79, 0x09, 0x2d, 0xfe, 0x20, 0xee, 0xd5,
	0xb3, 0xf6, 0x25, 0x65, 0x5a, 0xbb, 0x4c, 0x5b, 0xfc, 0x81, 0x3c, 0x01, 0x0d, 0x19, 0x4b, 0xd9,
	0x6c, 0x95, 0x45, 0xa5, 0x5f, 0x5d, 0x91, 0x38, 0xcf, 0x22, 0xe3, 0x3d, 0xc0, 0x97, 0x84, 0xfd,
	0xbb, 0x8c, 0x33, 0xe8, 0x5d, 0xc7, 0x51, 0x82, 0x0b, 0xe1, 0x22, 0x79, 0x0a, 0x5a, 0x16, 0x47,
	0x49, 0xc0, 0xef, 0x59, 0xe1, 0x73, 0x9f, 0xd6, 0x09, 0xf2, 0xbc, 0x1c, 0x83, 0xbd, 0xe1, 0x98,
	0x09, 0x09, 0x7d, 0xda, 0xc8, 0x18, 0xbf, 0x5a, 0xd0, 0x2e, 0xfa, 0x4c, 0xa0, 0x2b, 0xc5, 0x94,
	0xd7, 0xaa, 0x24, 0x48, 0xaf, 0x3c, 0x85, 0x56, 0x18, 0xf2, 0x0a, 0xda, 0xf3, 0xbb, 0x34, 0xbc,
	0x2d, 0x27, 0x34, 0x98, 0x94, 0x1b, 0x69, 0xe7, 0x49, 0x4f, 0xa1, 0x45, 0x95, 0x9c, 0xc0, 0x7f,
	0x7f, 0xec, 0xa5, 0x98, 0x4b, 0xcf, 0x3a, 0xf8, 0x6b, 0xa4, 0x42, 0x87, 0xa7, 0xd0, 0xdd, 0xf0,
	0x51, 0x86, 0xbc, 0x05, 0x8d, 0x49, 0xdf, 0xf5, 0x6d, 0x41, 0xde, 0xab, 0xa5, 0x95, 0x05, 0x4f,
	0xa1, 0x35, 0x8a, 0x1c, 0x03, 0xdc, 0x57, 0xde, 0xea, 0x6d, 0xc1, 0x21, 0x92, 0x53, 0xbb, 0xee,
	0x29, 0xb4, 0x81, 0x13, 0xbb, 0xc3, 0x30, 0xe0, 0x29, 0xd3, 0x3b, 0xc2, 0x29, 0x19, 0xda, 0x3b,
	0xa5, 0x4b, 0x6f, 0x6c, 0xd0, 0xaa, 0xe5, 0x25, 0x7d, 0xe8, 0x52, 0xf7, 0xa3, 0x7f, 0x3d, 0x75,
	0xe9, 0x50, 0x21, 0x1a, 0xb4, 0xed, 0xcf, 0x97, 0xce, 0xd9, 0x50, 0x25, 0x03, 0xd0, 0x1c, 0xef,
	0xc4, 0xbf, 0x70, 0x2e, 0x4f, 0xdd, 0x61, 0x2b, 0x0f, 0xa9, 0xfb, 0xc9, 0x75, 0xa6, 0xfe, 0xe5,
	0xc5, 0x70, 0xcb, 0xfa, 0x00, 0x1d, 0xd1, 0x23, 0x23, 0x47, 0xb0, 0xed, 0x2c, 0x03, 0x4e, 0xaa,
	0x05, 0x6a, 0x0c, 0x76, 0x34, 0x78, 0xf4, 0x5a, 0x0c, 0x65, 0xac, 0x1e, 0xa9, 0xf6, 0x37, 0x30,
	0x52, 0x16, 0x4d, 0x96, 0x9b, 0x35, 0xb2, 0x3b, 0x5c, 0x44, 0xc8, 0x26, 0x37, 0xc1, 0x9c, 0xc5,
	0xa1, 0x04, 0xe7, 0xaf, 0xdd, 0x1e, 0x14, 0xfd, 0xaf, 0x82, 0xf0, 0x36, 0x88, 0xf0, 0xeb, 0xeb,
	0x28, 0xe6, 0xcb, 0xfb, 0x79, 0x3e, 0x21, 0xb3, 0xc1, 0x34, 0x0b, 0xa6, 0x59, 0x30, 0xcd, 0x9c,
	0x39, 0x2f, 0xfe, 0x66, 0xde, 0xfd, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x19, 0x83, 0x34, 0xd3, 0x82,
	0x04, 0x00, 0x00,
}
