package main

import (
	"net"

	"encoding/json"

	"github.com/CossackPyra/pyraconv"
	api "github.com/krecu/go-visitor/protoc/visitor"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RpcService struct {
	server *grpc.Server
	listen net.Listener
	app    *App
}

func NewRpcService(app *App) (proto *RpcService, err error) {

	proto = &RpcService{
		app:    app,
		server: grpc.NewServer(),
	}

	api.RegisterGreeterServer(proto.server, proto)
	reflection.Register(proto.server)
	proto.listen, err = net.Listen("tcp", proto.app.config.GetString("app.server.grpc.listen"))
	if err != nil {
		return
	}

	return
}

func (s *RpcService) Start() {
	go func(s *RpcService) {
		Logger.Fatal(s.server.Serve(s.listen))
	}(s)
}

// получение
func (s *RpcService) Get(ctx context.Context, in *api.GetRequest) (*api.Reply, error) {

	reply := &api.Reply{}
	if values, err := s.app.visitor.Get(pyraconv.ToString(in.Id)); err == nil {

		if values == nil {
			reply.Status = "0001"
		} else {
			reply.Status = "0000"
		}

		if buf, err := json.Marshal(values); err == nil {
			reply.Body = string(buf)
		} else {
			reply.Status = err.Error()
		}

	} else {
		reply.Status = err.Error()
	}

	return reply, nil
}

// создание
func (s *RpcService) Post(ctx context.Context, in *api.PostRequest) (*api.Reply, error) {

	var extra map[string]interface{}

	reply := &api.Reply{}

	if in.Id == "" || in.Ua == "" || in.Ip == "" || net.ParseIP(in.Ip) == nil {
		reply.Status = "0002"
		return reply, nil
	}

	json.Unmarshal([]byte(in.Extra), &extra)

	if values, err := s.app.visitor.Post(in.Id, in.Ip, in.Ua, extra); err == nil {
		if buf, err := json.Marshal(values); err == nil {
			reply.Body = string(buf)
		} else {
			reply.Status = err.Error()
		}
		reply.Status = "0000"
	} else {

		if err == VisitorErrorEmpty {
			reply.Status = "0001"
		} else {
			reply.Status = err.Error()
		}
	}
	return reply, nil
}

// изменение
func (s *RpcService) Patch(ctx context.Context, in *api.PatchRequest) (*api.Reply, error) {
	return &api.Reply{Body: "", Status: ""}, nil
}

// удаление
func (s *RpcService) Delete(ctx context.Context, in *api.DeleteRequest) (*api.Reply, error) {

	reply := &api.Reply{}

	if err := s.app.visitor.Delete(pyraconv.ToString(in.Id)); err != nil {
		reply.Status = err.Error()
	} else {
		reply.Status = "0000"
	}

	return reply, nil
}
