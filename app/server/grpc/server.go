package grpc

import (
	"fmt"

	"time"

	"net"

	"github.com/CossackPyra/pyraconv"
	"github.com/krecu/go-visitor/app/module/cache"
	cacheAeroSpike "github.com/krecu/go-visitor/app/module/cache/aerospike"
	"github.com/krecu/go-visitor/app/module/provider/device"
	"github.com/krecu/go-visitor/app/module/provider/device/browscap"
	"github.com/krecu/go-visitor/app/module/provider/device/uasurfer"
	"github.com/krecu/go-visitor/app/module/provider/geo"
	"github.com/krecu/go-visitor/app/module/provider/geo/maxmind"
	"github.com/krecu/go-visitor/app/module/provider/geo/sypexgeo"
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
	Db                   cacheAeroSpike.Option
	Visitor              struct {
		SypexGeo sypexgeo.Option
		MaxMind  maxmind.Option
		BrowsCap browscap.Option
		UaSurfer uasurfer.Option
	}
}

type Server struct {
	opt     Options
	server  *grpc.Server
	visitor *wrapper.Wrapper
	db      cache.Cache
}

func New(opt Options) (proto *Server, err error) {

	var (
		sp, mm geo.Geo
		br, ua device.Device
	)

	proto = &Server{
		opt: opt,
	}

	// инстацируем сервер
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

	// подключаем БД
	proto.db, err = cacheAeroSpike.New(proto.opt.Db)
	if err != nil {
		err = fmt.Errorf("DB: %s", err)
		return
	}

	// формируем список провайдеров и добавляем в визитор
	proto.visitor = wrapper.New()

	// Основная GEO база для RU трафика
	if sp, err = sypexgeo.New(proto.opt.Visitor.SypexGeo); err == nil {
		proto.visitor.AddGeoProvider(sp)
	} else {
		err = fmt.Errorf("SYPEXGEO: %s", err)
		return
	}

	// Основная ГЕО база для не RU трафика
	if mm, err = maxmind.New(proto.opt.Visitor.MaxMind); err == nil {
		proto.visitor.AddGeoProvider(mm)
	} else {
		err = fmt.Errorf("MAXMIND: %s", err)
		return
	}

	// Основная база User-Agent
	if br, err = browscap.New(proto.opt.Visitor.BrowsCap); err == nil {
		proto.visitor.AddDeviceProvider(br)
	} else {
		err = fmt.Errorf("BROWSCAP: %s", err)
		return
	}

	// Словарное определение User-Agent
	if ua, err = uasurfer.New(proto.opt.Visitor.UaSurfer); err == nil {
		proto.visitor.AddDeviceProvider(ua)
	} else {
		err = fmt.Errorf("UASURFER: %s", err)
		return
	}

	return
}

func (s *Server) Start() (err error) {

	listen, err := net.Listen("tcp", s.opt.Listen)
	if err == nil {
		return s.server.Serve(listen)
	}

	return
}

func (s *Server) Get(ctx context.Context, request *greeter.GetRequest) (reply *greeter.Reply, err error) {

	var (
		data = make(map[string]interface{})
		info *wrapper.Model
		buf  []byte
	)

	reply = &greeter.Reply{}

	err = s.db.Get(request.GetId(), &data)
	if err == nil {

		info, err = s.visitor.Parse(pyraconv.ToString(data[FieldIpV4]), pyraconv.ToString(data[FieldPersonalUa]))
		if err != nil {
			reply.Status = fmt.Sprintf("PARSING: %s", err)
			return
		}

		model := &Model{}
		err = model.Formed(info, data)

		if err != nil {
			reply.Status = fmt.Sprintf("MODEL: Formed with err - %s", err)
			return
		}

		buf, err = model.Marshal()
		if err != nil {
			reply.Status = fmt.Sprintf("MODEL: Marshal with err - %s", err)
			return
		}

		reply.Body = string(buf)
	}

	return
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

	info, err := s.visitor.Parse(request.GetIp(), request.GetUa())
	if err != nil {
		reply.Status = fmt.Sprintf("PARSING: %s", err)
		return
	}

	model := &Model{}
	buf, err = model.UnMarshal(info, map[string]interface{}{
		FieldId:         request.GetId(),
		FieldCreated:    time.Now().Unix(),
		FieldUpdated:    time.Now().Unix(),
		FieldPersonalUa: request.GetUa(),
		FieldIpV4:       request.GetIp(),
		FieldExtra:      request.GetExtra(),
	})

	if err != nil {
		reply.Status = fmt.Sprintf("JSON: %s", err)
		return
	}

	reply.Body = string(buf)

	// сохраняем данные в дб
	go func() {

		err = s.db.Set(request.GetId(), map[string]interface{}{
			FieldId:         request.GetId(),
			FieldCreated:    time.Now().Unix(),
			FieldUpdated:    time.Now().Unix(),
			FieldPersonalUa: request.GetUa(),
			FieldIpV4:       request.GetIp(),
			FieldExtra:      request.GetExtra(),
		})
		if err != nil {
			fmt.Println(err)
		}
	}()

	return
}
func (s *Server) Patch(context.Context, *greeter.PatchRequest) (*greeter.Reply, error) {
	reply := &greeter.Reply{}
	return reply, nil
}
