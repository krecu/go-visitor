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
)

type server struct{
	Core *visitor.Visitor
}

// основной метод получения данных пользователя
func (s *server) GetVisitor(ctx context.Context, in *pb.VisitorRequest) (*pb.VisitorReply, error) {

	info, err := s.Core.Identify(in.GetIp(), in.GetUa()); if err != nil {
		//in.GetId()
	}

	// упаковываем структуру в json
	jsonData, err := json.Marshal(info)

	if err != nil {
		return &pb.VisitorReply{Status: "false", Body: err.Error()}, nil
	}

	return &pb.VisitorReply{Status: "ok", Body: string(jsonData)}, nil
}

func main()  {

	cpu := flag.Int("cpu", 1, "Count usage cpu")
	addr := flag.String("addr", ":50051", "Network")
	db := flag.String("db", "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor", "Path to database")
	buffer := flag.Int("buffer", 100000, "Путь к конфигам и базам")
	debug := flag.Bool("debug", true, "Путь к конфигам и базам")
	flag.Parse()

	runtime.GOMAXPROCS(int(*cpu))

	core, err := visitor.New(*debug, *db, *buffer); if err != nil {
		panic(err)
	}

	// вешаем листнера на порт
	listen, err := net.Listen("tcp", *addr); if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{
		Core: core,
	})
	reflection.Register(s)
	if err := s.Serve(listen); err != nil {
	}
}