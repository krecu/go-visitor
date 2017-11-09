// Code generated by protoc-gen-go.
// source: visitor.proto
// DO NOT EDIT!

/*
Package rpc is a generated protocol buffer package.

It is generated from these files:
	visitor.proto

It has these top-level messages:
	GetRequest
	DeleteRequest
	PostRequest
	PatchRequest
	Reply
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
type GetRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// The request message containing the user's name.
type DeleteRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DeleteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// The request message containing the user's name.
type PostRequest struct {
	Ip    string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Ua    string `protobuf:"bytes,2,opt,name=ua" json:"ua,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Extra string `protobuf:"bytes,4,opt,name=extra" json:"extra,omitempty"`
}

func (m *PostRequest) Reset()                    { *m = PostRequest{} }
func (m *PostRequest) String() string            { return proto.CompactTextString(m) }
func (*PostRequest) ProtoMessage()               {}
func (*PostRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PostRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *PostRequest) GetUa() string {
	if m != nil {
		return m.Ua
	}
	return ""
}

func (m *PostRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PostRequest) GetExtra() string {
	if m != nil {
		return m.Extra
	}
	return ""
}

// The request message containing the user's name.
type PatchRequest struct {
	Ip    string `protobuf:"bytes,1,opt,name=ip" json:"ip,omitempty"`
	Ua    string `protobuf:"bytes,2,opt,name=ua" json:"ua,omitempty"`
	Id    string `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
	Extra string `protobuf:"bytes,4,opt,name=extra" json:"extra,omitempty"`
}

func (m *PatchRequest) Reset()                    { *m = PatchRequest{} }
func (m *PatchRequest) String() string            { return proto.CompactTextString(m) }
func (*PatchRequest) ProtoMessage()               {}
func (*PatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PatchRequest) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *PatchRequest) GetUa() string {
	if m != nil {
		return m.Ua
	}
	return ""
}

func (m *PatchRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PatchRequest) GetExtra() string {
	if m != nil {
		return m.Extra
	}
	return ""
}

// The response message containing the greetings
type Reply struct {
	Status string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Body   string `protobuf:"bytes,2,opt,name=body" json:"body,omitempty"`
}

func (m *Reply) Reset()                    { *m = Reply{} }
func (m *Reply) String() string            { return proto.CompactTextString(m) }
func (*Reply) ProtoMessage()               {}
func (*Reply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Reply) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Reply) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*GetRequest)(nil), "rpc.GetRequest")
	proto.RegisterType((*DeleteRequest)(nil), "rpc.DeleteRequest")
	proto.RegisterType((*PostRequest)(nil), "rpc.PostRequest")
	proto.RegisterType((*PatchRequest)(nil), "rpc.PatchRequest")
	proto.RegisterType((*Reply)(nil), "rpc.Reply")
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
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Reply, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Reply, error)
	Post(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*Reply, error)
	Patch(ctx context.Context, in *PatchRequest, opts ...grpc.CallOption) (*Reply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/rpc.Greeter/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/rpc.Greeter/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) Post(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/rpc.Greeter/Post", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterClient) Patch(ctx context.Context, in *PatchRequest, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/rpc.Greeter/Patch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterServer interface {
	// Sends a greeting
	Get(context.Context, *GetRequest) (*Reply, error)
	Delete(context.Context, *DeleteRequest) (*Reply, error)
	Post(context.Context, *PostRequest) (*Reply, error)
	Patch(context.Context, *PatchRequest) (*Reply, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_Post_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).Post(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/Post",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).Post(ctx, req.(*PostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Greeter_Patch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).Patch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Greeter/Patch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).Patch(ctx, req.(*PatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Greeter_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Greeter_Delete_Handler,
		},
		{
			MethodName: "Post",
			Handler:    _Greeter_Post_Handler,
		},
		{
			MethodName: "Patch",
			Handler:    _Greeter_Patch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "visitor.proto",
}

func init() { proto.RegisterFile("visitor.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x91, 0x4f, 0x4b, 0xc3, 0x40,
	0x10, 0xc5, 0xcd, 0x5f, 0x71, 0x6c, 0xad, 0x0e, 0x22, 0x41, 0x04, 0x35, 0x88, 0x14, 0x0f, 0x39,
	0xd8, 0x6f, 0x50, 0x84, 0x5c, 0x43, 0x14, 0xef, 0x69, 0x32, 0xe8, 0x42, 0x61, 0xd7, 0xdd, 0x89,
	0xd8, 0x0f, 0xe5, 0x77, 0x94, 0xec, 0xae, 0xda, 0x14, 0xbc, 0x79, 0xdb, 0xf7, 0xf6, 0xcd, 0x3b,
	0xfc, 0x1e, 0x4c, 0xdf, 0x85, 0x11, 0x2c, 0x75, 0xa1, 0xb4, 0x64, 0x89, 0x91, 0x56, 0x6d, 0x7e,
	0x01, 0x50, 0x12, 0xd7, 0xf4, 0xd6, 0x93, 0x61, 0x3c, 0x82, 0x50, 0x74, 0x59, 0x70, 0x15, 0xcc,
	0x0f, 0xea, 0x50, 0x74, 0xf9, 0x25, 0x4c, 0x1f, 0x68, 0x4d, 0x4c, 0x7f, 0x05, 0x1e, 0xe1, 0xb0,
	0x92, 0x66, 0x74, 0xaf, 0x7e, 0xbe, 0xd5, 0xa0, 0xfb, 0x26, 0x0b, 0x9d, 0xee, 0x1b, 0x7f, 0x1e,
	0x7d, 0x9f, 0xe3, 0x29, 0x24, 0xf4, 0xc1, 0xba, 0xc9, 0x62, 0x6b, 0x39, 0x91, 0x3f, 0xc1, 0xa4,
	0x6a, 0xb8, 0x7d, 0xfd, 0xdf, 0xd6, 0x05, 0x24, 0x35, 0xa9, 0xf5, 0x06, 0xcf, 0x20, 0x35, 0xdc,
	0x70, 0x6f, 0x7c, 0xa5, 0x57, 0x88, 0x10, 0xaf, 0x64, 0xb7, 0xf1, 0xc5, 0xf6, 0x7d, 0xff, 0x19,
	0xc0, 0x7e, 0xa9, 0x89, 0x98, 0x34, 0xde, 0x40, 0x54, 0x12, 0xe3, 0xac, 0xd0, 0xaa, 0x2d, 0x7e,
	0xa1, 0x9d, 0x83, 0x35, 0x6c, 0x77, 0xbe, 0x87, 0x77, 0x90, 0x3a, 0x64, 0x88, 0xd6, 0x1f, 0xf1,
	0xdb, 0xc9, 0xde, 0x42, 0x3c, 0xd0, 0xc3, 0x63, 0xeb, 0x6e, 0x81, 0xdc, 0xc9, 0xcd, 0x21, 0xb1,
	0x40, 0xf0, 0xc4, 0x05, 0xb7, 0xe0, 0x8c, 0x93, 0xcb, 0x6b, 0x98, 0x09, 0x59, 0xbc, 0x0c, 0x96,
	0x1f, 0x7b, 0x39, 0x79, 0x76, 0x8f, 0x6a, 0x18, 0xbd, 0x0a, 0x56, 0xa9, 0x5d, 0x7f, 0xf1, 0x15,
	0x00, 0x00, 0xff, 0xff, 0x6e, 0xf6, 0x1d, 0x0a, 0x0e, 0x02, 0x00, 0x00,
}