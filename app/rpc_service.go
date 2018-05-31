package main

import (
	"net"

	"encoding/json"

	"time"

	"github.com/CossackPyra/pyraconv"
	api "github.com/krecu/go-visitor/protoc/visitor"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type RpcService struct {
	server *grpc.Server
	listen net.Listener
	app    *App
}

func NewRpcService(app *App) (proto *RpcService, err error) {

	proto = &RpcService{
		app: app,
		server: grpc.NewServer(
			grpc.MaxConcurrentStreams(100),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: 0,
				//MaxConnectionAge:      100 * time.Millisecond,
				//MaxConnectionAgeGrace: 100 * time.Millisecond,
				Time:    1 * time.Minute,
				Timeout: 1 * time.Second,
			}),
		),
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

	var _total time.Duration

	reply := &api.Reply{}

	if values, err := s.app.visitor.Get(pyraconv.ToString(in.Id)); err == nil {

		_total = values.Debug.TimeTotal

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

	Logger.WithFields(logrus.Fields{
		"state":    reply.Status,
		"op":       "Get",
		"duration": _total.Seconds() * 1000,
	}).Debugf("GET")

	return reply, nil
}

// создание
func (s *RpcService) Post(ctx context.Context, in *api.PostRequest) (*api.Reply, error) {

	var extra map[string]interface{}
	var _total time.Duration

	reply := &api.Reply{}

	if in.Id == "" || in.Ua == "" || in.Ip == "" || net.ParseIP(in.Ip) == nil {
		reply.Status = "0002"
		return reply, nil
	}

	json.Unmarshal([]byte(in.Extra), &extra)

	if values, err := s.app.visitor.Post(in.Id, in.Ip, in.Ua, extra); err == nil {

		_total = values.Debug.TimeTotal

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

	Logger.WithFields(logrus.Fields{
		"state":    reply.Status,
		"op":       "Post",
		"duration": _total.Seconds() * 1000,
	}).Debugf("POST")

	return reply, nil
}

// изменение
func (s *RpcService) Patch(ctx context.Context, in *api.PatchRequest) (*api.Reply, error) {

	//_total := time.Now()

	var fields map[string]interface{}

	reply := &api.Reply{}

	if in.Id == "" || in.Fields == "" {
		reply.Status = "0002"
		return reply, nil
	}

	json.Unmarshal([]byte(in.Fields), &fields)

	if values, err := s.app.visitor.Patch(in.Id, fields); err == nil {
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

	//Logger.WithFields(logrus.Fields{
	//	"op":       "Patch",
	//	"duration": time.Since(_total).Seconds(),
	//}).Debugf("Patch: %f", time.Since(_total).Seconds())

	return reply, nil
}

// удаление
func (s *RpcService) Delete(ctx context.Context, in *api.DeleteRequest) (*api.Reply, error) {

	_total := time.Now()

	reply := &api.Reply{}

	if err := s.app.visitor.Delete(pyraconv.ToString(in.Id)); err != nil {
		reply.Status = err.Error()
	} else {
		reply.Status = "0000"
	}

	Logger.WithFields(logrus.Fields{
		"op":       "Delete",
		"duration": time.Since(_total).Seconds(),
	}).Debugf("Delete: %f", time.Since(_total).Seconds())

	return reply, nil
}
