// Code generated by protoc-gen-go.
// source: visitor.proto
// DO NOT EDIT!

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	visitor.proto

It has these top-level messages:
	VisitorRequest
	VisitorReply
*/
package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type VisitorRequest struct {
	Ip    string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Ua    string `protobuf:"bytes,2,opt,name=ua" json:"ua,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Extra string `protobuf:"bytes,4,opt,name=extra" json:"extra,omitempty"`
}

func (m *VisitorRequest) Reset()                    { *m = VisitorRequest{} }
func (m *VisitorRequest) String() string            { return proto.CompactTextString(m) }
func (*VisitorRequest) ProtoMessage()               {}
func (*VisitorRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *VisitorRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *VisitorRequest) GetUa() string {
	if m != nil {
		return m.Ua
	}
	return ""
}

func (m *VisitorRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *VisitorRequest) GetExtra() string {
	if m != nil {
		return m.Extra
	}
	return ""
}

// The response message containing the greetings
type VisitorReply struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Body   string `protobuf:"bytes,2,opt,name=body" json:"body,omitempty"`
}

func (m *VisitorReply) Reset()                    { *m = VisitorReply{} }
func (m *VisitorReply) String() string            { return proto.CompactTextString(m) }
func (*VisitorReply) ProtoMessage()               {}
func (*VisitorReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *VisitorReply) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *VisitorReply) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*VisitorRequest)(nil), "rpc.VisitorRequest")
	proto.RegisterType((*VisitorReply)(nil), "rpc.VisitorReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Greeter service

type GreeterClient interface {
	// Sends a greeting
	GetVisitor(ctx context.Context, in *VisitorRequest, opts ...grpc.CallOption) (*VisitorReply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) GetVisitor(ctx context.Context, in *VisitorRequest, opts ...grpc.CallOption) (*VisitorReply, error) {
	out := new(VisitorReply)
	err := grpc.Invoke(ctx, "/rpc.Greeter/GetVisitor", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) PutVisitor(ctx context.Context, in *VisitorRequest, opts ...grpc.CallOption) (*VisitorReply, error) {
	out := new(VisitorReply)
	err := grpc.Invoke(ctx, "/rpc.Greeter/PutVisitor", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterServer interface {
	// Sends a greeting
	GetVisitor(context.Context, *VisitorRequest) (*VisitorReply, error)
	PutVisitor(context.Context, *VisitorRequest) (*VisitorReply, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_GetVisitor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VisitorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).GetVisitor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/GetVisitor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).GetVisitor(ctx, req.(*VisitorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_PutVisitor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VisitorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).PutVisitor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/PutVisitor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).PutVisitor(ctx, req.(*VisitorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVisitor",
			Handler:    _Greeter_GetVisitor_Handler,
		},
		{
			MethodName: "PutVisitor",
			Handler:    _Greeter_PutVisitor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "visitor.proto",
}

func init() { proto.RegisterFile("visitor.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 221 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x31, 0x4b, 0x03, 0x41,
	0x10, 0x85, 0xbd, 0x4b, 0x8c, 0x38, 0x68, 0xc4, 0x51, 0xe4, 0xd0, 0x46, 0xae, 0xb2, 0x5a, 0x44,
	0xc1, 0xc2, 0xce, 0x34, 0xb1, 0x0c, 0x29, 0x62, 0xbd, 0xc9, 0x0d, 0xba, 0xb0, 0xb2, 0xe3, 0xec,
	0x9c, 0xe6, 0xfe, 0xbd, 0xec, 0xdd, 0x12, 0xb8, 0xee, 0x7d, 0xef, 0xc1, 0xc7, 0x30, 0x70, 0xfe,
	0xeb, 0xa2, 0xd3, 0x20, 0x86, 0x25, 0x68, 0xc0, 0x89, 0xf0, 0xae, 0xde, 0xc0, 0x7c, 0x33, 0xb4,
	0x6b, 0xfa, 0x69, 0x29, 0x2a, 0xce, 0xa1, 0x74, 0x5c, 0x15, 0xf7, 0xc5, 0xc3, 0xe9, 0xba, 0x74,
	0x9c, 0xb8, 0xb5, 0x55, 0x39, 0x70, 0x6b, 0xfb, 0xbd, 0xa9, 0x26, 0x79, 0x6f, 0xf0, 0x1a, 0x8e,
	0x69, 0xaf, 0x62, 0xab, 0x69, 0x5f, 0x0d, 0x50, 0xbf, 0xc2, 0xd9, 0xc1, 0xcb, 0xbe, 0xc3, 0x1b,
	0x98, 0x45, 0xb5, 0xda, 0xc6, 0x6c, 0xce, 0x84, 0x08, 0xd3, 0x6d, 0x68, 0xba, 0xec, 0xef, 0xf3,
	0xd3, 0x1b, 0x9c, 0x2c, 0x85, 0x48, 0x49, 0xf0, 0x05, 0x60, 0x49, 0x9a, 0x4d, 0x78, 0x65, 0x84,
	0x77, 0x66, 0x7c, 0xef, 0xed, 0xe5, 0xb8, 0x64, 0xdf, 0xd5, 0x47, 0x8b, 0x47, 0xb8, 0x73, 0xc1,
	0x7c, 0xa6, 0x85, 0xf6, 0xf6, 0x9b, 0x3d, 0x45, 0xf3, 0x45, 0xde, 0x87, 0xbf, 0x20, 0xbe, 0x59,
	0x5c, 0xbc, 0xa7, 0xfc, 0x91, 0xf2, 0x2a, 0xfd, 0x62, 0x55, 0x6c, 0x67, 0xfd, 0x53, 0x9e, 0xff,
	0x03, 0x00, 0x00, 0xff, 0xff, 0x89, 0x95, 0xa3, 0x6c, 0x25, 0x01, 0x00, 0x00,
}
