package grpc

import (
	"encoding/json"

	"fmt"

	"time"

	"net"

	"github.com/krecu/go-visitor/app/module/wrapper"
	greeter "github.com/krecu/go-visitor/app/protoc/visitor"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Options struct {
	MaxConcurrentStreams uint32
	MaxConnectionIdle    time.Duration
	Time                 time.Duration
	Timeout              time.Duration
	Listen               string
	visitor              *wrapper.Wrapper
}

type Server struct {
	opt    Options
	server *grpc.Server
}

func New(opt Options) (proto *Server, err error) {

	proto = &Server{
		opt: opt,
	}

	proto.server = grpc.NewServer(
		grpc.MaxConcurrentStreams(proto.opt.MaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: proto.opt.MaxConnectionIdle * time.Second,
			Time:              proto.opt.Time * time.Second,
			Timeout:           proto.opt.Timeout * time.Second,
		}),
	)

	greeter.RegisterGreeterServer(proto.server, proto)
	reflection.Register(proto.server)

	return
}

func (s *Server) Start() (err error) {

	listen, err := net.Listen("tcp", s.opt.Listen)
	if err == nil {
		return s.server.Serve(listen)
	}

	return
}

func (s *Server) Get(ctx context.Context, request *greeter.GetRequest) (*greeter.Reply, error) {
	reply := &greeter.Reply{}
	return reply, nil
}
func (s *Server) Delete(ctx context.Context, request *greeter.DeleteRequest) (*greeter.Reply, error) {
	reply := &greeter.Reply{}
	return reply, nil
}
func (s *Server) Post(ctx context.Context, request *greeter.PostRequest) (reply *greeter.Reply, err error) {

	var (
		buf []byte
	)

	reply = &greeter.Reply{}

	data, err := s.visitor.Parse(request.GetIp(), request.GetUa())
	if err != nil {
		reply.Status = fmt.Sprintf("PARSING: %s", err)
		return
	}

	buf, err = json.Marshal(data)
	if err != nil {
		reply.Status = fmt.Sprintf("JSON: %s", err)
		return
	}

	reply.Body = string(buf)

	return
}
func (s *Server) Patch(context.Context, *greeter.PatchRequest) (*greeter.Reply, error) {
	reply := &greeter.Reply{}
	return reply, nil
}
