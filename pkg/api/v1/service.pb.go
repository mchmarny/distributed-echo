// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package echo

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type BroadcastMessage struct {
	Self                 *Node    `protobuf:"bytes,1,opt,name=self,proto3" json:"self,omitempty"`
	Targets              []*Node  `protobuf:"bytes,2,rep,name=targets,proto3" json:"targets,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BroadcastMessage) Reset()         { *m = BroadcastMessage{} }
func (m *BroadcastMessage) String() string { return proto.CompactTextString(m) }
func (*BroadcastMessage) ProtoMessage()    {}
func (*BroadcastMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *BroadcastMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BroadcastMessage.Unmarshal(m, b)
}
func (m *BroadcastMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BroadcastMessage.Marshal(b, m, deterministic)
}
func (m *BroadcastMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BroadcastMessage.Merge(m, src)
}
func (m *BroadcastMessage) XXX_Size() int {
	return xxx_messageInfo_BroadcastMessage.Size(m)
}
func (m *BroadcastMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_BroadcastMessage.DiscardUnknown(m)
}

var xxx_messageInfo_BroadcastMessage proto.InternalMessageInfo

func (m *BroadcastMessage) GetSelf() *Node {
	if m != nil {
		return m.Self
	}
	return nil
}

func (m *BroadcastMessage) GetTargets() []*Node {
	if m != nil {
		return m.Targets
	}
	return nil
}

type BroadcastResult struct {
	Count                int32    `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BroadcastResult) Reset()         { *m = BroadcastResult{} }
func (m *BroadcastResult) String() string { return proto.CompactTextString(m) }
func (*BroadcastResult) ProtoMessage()    {}
func (*BroadcastResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *BroadcastResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BroadcastResult.Unmarshal(m, b)
}
func (m *BroadcastResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BroadcastResult.Marshal(b, m, deterministic)
}
func (m *BroadcastResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BroadcastResult.Merge(m, src)
}
func (m *BroadcastResult) XXX_Size() int {
	return xxx_messageInfo_BroadcastResult.Size(m)
}
func (m *BroadcastResult) XXX_DiscardUnknown() {
	xxx_messageInfo_BroadcastResult.DiscardUnknown(m)
}

var xxx_messageInfo_BroadcastResult proto.InternalMessageInfo

func (m *BroadcastResult) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type Node struct {
	Region               string   `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	Uri                  string   `protobuf:"bytes,2,opt,name=uri,proto3" json:"uri,omitempty"`
	Port                 string   `protobuf:"bytes,3,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{2}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetRegion() string {
	if m != nil {
		return m.Region
	}
	return ""
}

func (m *Node) GetUri() string {
	if m != nil {
		return m.Uri
	}
	return ""
}

func (m *Node) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

type EchoMessage struct {
	Id                   string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Sent                 *timestamp.Timestamp `protobuf:"bytes,2,opt,name=sent,proto3" json:"sent,omitempty"`
	Source               *Node                `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	Target               *Node                `protobuf:"bytes,4,opt,name=target,proto3" json:"target,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *EchoMessage) Reset()         { *m = EchoMessage{} }
func (m *EchoMessage) String() string { return proto.CompactTextString(m) }
func (*EchoMessage) ProtoMessage()    {}
func (*EchoMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{3}
}

func (m *EchoMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoMessage.Unmarshal(m, b)
}
func (m *EchoMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoMessage.Marshal(b, m, deterministic)
}
func (m *EchoMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoMessage.Merge(m, src)
}
func (m *EchoMessage) XXX_Size() int {
	return xxx_messageInfo_EchoMessage.Size(m)
}
func (m *EchoMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoMessage.DiscardUnknown(m)
}

var xxx_messageInfo_EchoMessage proto.InternalMessageInfo

func (m *EchoMessage) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *EchoMessage) GetSent() *timestamp.Timestamp {
	if m != nil {
		return m.Sent
	}
	return nil
}

func (m *EchoMessage) GetSource() *Node {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *EchoMessage) GetTarget() *Node {
	if m != nil {
		return m.Target
	}
	return nil
}

func init() {
	proto.RegisterType((*BroadcastMessage)(nil), "echo.BroadcastMessage")
	proto.RegisterType((*BroadcastResult)(nil), "echo.BroadcastResult")
	proto.RegisterType((*Node)(nil), "echo.Node")
	proto.RegisterType((*EchoMessage)(nil), "echo.EchoMessage")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x51, 0xc1, 0x4e, 0xc2, 0x40,
	0x10, 0x85, 0x52, 0x30, 0x0c, 0x51, 0x71, 0xa2, 0xa4, 0xe9, 0x41, 0xc9, 0xc6, 0x44, 0x4f, 0x4b,
	0x82, 0x57, 0x4f, 0x46, 0x8f, 0x7a, 0xa8, 0x1e, 0xbc, 0x96, 0x76, 0x28, 0x9b, 0x00, 0x43, 0x76,
	0xb7, 0x5e, 0xfc, 0x0e, 0xff, 0xd7, 0xec, 0x2e, 0x25, 0x5a, 0x6e, 0x33, 0xef, 0xbd, 0xbe, 0xee,
	0x7b, 0x03, 0xa7, 0x86, 0xf4, 0x97, 0x2a, 0x48, 0xee, 0x34, 0x5b, 0xc6, 0x98, 0x8a, 0x15, 0xa7,
	0x37, 0x15, 0x73, 0xb5, 0xa6, 0x99, 0xc7, 0x16, 0xf5, 0x72, 0x66, 0xd5, 0x86, 0x8c, 0xcd, 0x37,
	0xbb, 0x20, 0x13, 0x9f, 0x30, 0x7e, 0xd2, 0x9c, 0x97, 0x45, 0x6e, 0xec, 0x2b, 0x19, 0x93, 0x57,
	0x84, 0xd7, 0x10, 0x1b, 0x5a, 0x2f, 0x93, 0xee, 0xb4, 0x7b, 0x3f, 0x9a, 0x83, 0x74, 0x4e, 0xf2,
	0x8d, 0x4b, 0xca, 0x3c, 0x8e, 0xb7, 0x70, 0x62, 0x73, 0x5d, 0x91, 0x35, 0x49, 0x34, 0xed, 0xb5,
	0x24, 0x0d, 0x25, 0xee, 0xe0, 0xfc, 0xe0, 0x9c, 0x91, 0xa9, 0xd7, 0x16, 0x2f, 0xa1, 0x5f, 0x70,
	0xbd, 0xb5, 0xde, 0xb9, 0x9f, 0x85, 0x45, 0x3c, 0x43, 0xec, 0xbe, 0xc4, 0x09, 0x0c, 0x34, 0x55,
	0x8a, 0xb7, 0x9e, 0x1e, 0x66, 0xfb, 0x0d, 0xc7, 0xd0, 0xab, 0xb5, 0x4a, 0x22, 0x0f, 0xba, 0x11,
	0x11, 0xe2, 0x1d, 0x6b, 0x9b, 0xf4, 0x3c, 0xe4, 0x67, 0xf1, 0xd3, 0x85, 0xd1, 0x4b, 0xb1, 0xe2,
	0x26, 0xc4, 0x19, 0x44, 0xaa, 0xdc, 0x3b, 0x45, 0xaa, 0x44, 0xe9, 0x42, 0x6d, 0xad, 0xb7, 0x19,
	0xcd, 0x53, 0x19, 0x8a, 0x91, 0x4d, 0x31, 0xf2, 0xa3, 0x29, 0x26, 0xf3, 0x3a, 0x14, 0x30, 0x30,
	0x5c, 0xeb, 0x82, 0xfc, 0x5f, 0xfe, 0x67, 0xdc, 0x33, 0x4e, 0x13, 0xd2, 0x26, 0xf1, 0xb1, 0x26,
	0x30, 0xf3, 0xef, 0xf0, 0xac, 0xf7, 0x70, 0x1c, 0x7c, 0x84, 0xe1, 0xa1, 0x15, 0x9c, 0x04, 0x7d,
	0xfb, 0x00, 0xe9, 0x55, 0x0b, 0x0f, 0xf5, 0x89, 0x8e, 0x0b, 0xe1, 0xcc, 0xf0, 0x22, 0x08, 0xfe,
	0xe4, 0x4d, 0x8f, 0x21, 0xd1, 0x59, 0x0c, 0x7c, 0xbc, 0x87, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x0b, 0xd9, 0xe0, 0xa4, 0x1c, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EchoServiceClient is the client API for EchoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EchoServiceClient interface {
	Broadcast(ctx context.Context, in *BroadcastMessage, opts ...grpc.CallOption) (*BroadcastResult, error)
	Echo(ctx context.Context, in *EchoMessage, opts ...grpc.CallOption) (*EchoMessage, error)
}

type echoServiceClient struct {
	cc *grpc.ClientConn
}

func NewEchoServiceClient(cc *grpc.ClientConn) EchoServiceClient {
	return &echoServiceClient{cc}
}

func (c *echoServiceClient) Broadcast(ctx context.Context, in *BroadcastMessage, opts ...grpc.CallOption) (*BroadcastResult, error) {
	out := new(BroadcastResult)
	err := c.cc.Invoke(ctx, "/echo.EchoService/Broadcast", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoServiceClient) Echo(ctx context.Context, in *EchoMessage, opts ...grpc.CallOption) (*EchoMessage, error) {
	out := new(EchoMessage)
	err := c.cc.Invoke(ctx, "/echo.EchoService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServiceServer is the server API for EchoService service.
type EchoServiceServer interface {
	Broadcast(context.Context, *BroadcastMessage) (*BroadcastResult, error)
	Echo(context.Context, *EchoMessage) (*EchoMessage, error)
}

func RegisterEchoServiceServer(s *grpc.Server, srv EchoServiceServer) {
	s.RegisterService(&_EchoService_serviceDesc, srv)
}

func _EchoService_Broadcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BroadcastMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Broadcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.EchoService/Broadcast",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Broadcast(ctx, req.(*BroadcastMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.EchoService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServiceServer).Echo(ctx, req.(*EchoMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _EchoService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "echo.EchoService",
	HandlerType: (*EchoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Broadcast",
			Handler:    _EchoService_Broadcast_Handler,
		},
		{
			MethodName: "Echo",
			Handler:    _EchoService_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}