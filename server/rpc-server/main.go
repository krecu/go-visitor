package main

import (
	"runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/krecu/go-visitor/rpc"
	"google.golang.org/grpc/reflection"
	"net"
	"flag"
	"encoding/json"
	"github.com/krecu/go-visitor"
	"github.com/krecu/go-visitor/cache"
	"github.com/krecu/go-visitor/server"
	"log"
	"gopkg.in/natefinch/lumberjack.v2"
	"errors"
	"github.com/robertkowalski/graylog-golang"
	"os"
	"time"
	"fmt"
)

var (
	Conf *server.Config
	Core *server.Core
)

type RpcServer struct{
	Core *server.Core
}

// основной метод получения данных пользователя
func (s *RpcServer) GetVisitor(ctx context.Context, in *pb.VisitorRequest) (*pb.VisitorReply, error) {

	var err error
	var extra map[string]interface{}
	var jsonData []byte
	total := time.Now()

	if in.Id == "" || in.Ua == "" || in.Ip == "" {
		err = errors.New("Bad request")
	} else {

		if in.Extra != "" {
			err = json.Unmarshal([]byte(in.GetExtra()), &extra)
		}

		if Conf.Refresh {
			info, err := Core.Refresh(in.GetId(), in.GetIp(), in.GetUa(), extra);
			if err == nil {
				jsonData, err = json.Marshal(info)
			}
		} else {
			info, err := Core.Get(in.GetId(), in.GetIp(), in.GetUa(), extra);
			if err == nil {

				jsonData, err = json.Marshal(info)
			}
		}
	}

	if err != nil {
		LogNotify(LogMessage{
			ShortMessage: "Error: " + in.Id + ", " + in.Ua + ", " + in.Ip,
			State: "error",
		})
		return &pb.VisitorReply{Status: "false", Body: err.Error()}, nil
	} else {
		LogNotify(LogMessage{
			ShortMessage: string(jsonData),
			State: "ok",
			Duration: time.Now().Sub(total).Seconds(),
		})
	}

	return &pb.VisitorReply{Status: "ok", Body: string(jsonData)}, nil
}

// основной метод получения данных пользователя
func (s *RpcServer) PutVisitor(ctx context.Context, in *pb.VisitorRequest) (*pb.VisitorReply, error) {

	var err error
	var extra map[string]interface{}
	var jsonData []byte

	total := time.Now()

	if in.Id == "" || in.Ua == "" || in.Ip == "" {
		err = errors.New("Bad request")
	} else {

		if in.Extra != "" {
			err = json.Unmarshal([]byte(in.GetExtra()), &extra)
		}

		if Conf.Refresh {
			info, err := Core.Refresh(in.GetId(), in.GetIp(), in.GetUa(), extra);
			if err == nil {
				jsonData, err = json.Marshal(info)
			}
		} else {
			info, err := Core.Put(in.GetId(), in.GetIp(), in.GetUa(), extra);
			if err == nil {
				jsonData, err = json.Marshal(info)
			}
		}
	}

	if err != nil {
		LogNotify(LogMessage{
			ShortMessage: "Error: " + in.Id + ", " + in.Ua + ", " + in.Ip,
			State: "error",
		})
		return &pb.VisitorReply{Status: "false", Body: err.Error()}, nil
	} else {
		LogNotify(LogMessage{
			ShortMessage: string(jsonData),
			State: "ok",
			Duration: time.Now().Sub(total).Seconds(),
		})
	}

	return &pb.VisitorReply{Status: "ok", Body: string(jsonData)}, nil
}

type LogMessage struct {
	Version string		`json:"version"`
	Host string		`json:"host"`
	Timestamp int64		`json:"timestamp"`
	Facility string		`json:"facility"`
	ShortMessage string	`json:"short_message"`
	State string		`json:"state"`
	Duration float64	`json:"duration"`
}

func LogNotify(m LogMessage) {

	if Conf.Logger.Enabled {
		g := gelf.New(gelf.Config{
			GraylogPort:     Conf.Logger.Port,
			GraylogHostname: Conf.Logger.Host,
		})

		m.Host, _ = os.Hostname()
		m.Version = "1.0"
		m.Timestamp = time.Now().Unix()
		m.Facility = "Visitor"

		message, _ := json.Marshal(m)
		g.Log(string(message))
	}
}

func main()  {

	conf := flag.String("config", "./config.json", "Config file")

	flag.Parse()

	Conf = server.NewConfig(*conf)

	runtime.GOMAXPROCS(Conf.Cpu)

	if !Conf.Debug {
		log.SetOutput(&lumberjack.Logger{
			Filename:   Conf.Log,
			MaxSize:    500,
			MaxBackups: 3,
		})
	}

	cacheProvider := cache.New(Conf.AeroSpike.Host, Conf.AeroSpike.Port, Conf.AeroSpike.Ns, Conf.AeroSpike.Db, Conf.AeroSpike.Timeout, Conf.AeroSpike.Ttl)
	defer cacheProvider.Close()

	coreVisitor, err := visitor.New(Conf.Debug, Conf.Db, Conf.Buffer); if err != nil {
		panic(err)
	}

	Core = server.New(coreVisitor, cacheProvider)

	// вешаем листнера на порт
	listen, err := net.Listen("tcp", Conf.Listener); if err != nil {
		fmt.Println(err)
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &RpcServer{
		Core: Core,
	})
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
		fmt.Println(err)
		panic(err)
	}
}