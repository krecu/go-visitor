package main

import (
	"math/rand"
	"net"
	"runtime"

	"time"

	"bufio"
	"fmt"
	"os"

	"flag"

	"github.com/krecu/go-visitor/app/module/provider/device/browscap"
	"github.com/krecu/go-visitor/app/module/provider/device/uasurfer"
	"github.com/krecu/go-visitor/app/module/provider/geo/maxmind"
	"github.com/krecu/go-visitor/app/module/provider/geo/sypexgeo"
	"github.com/krecu/go-visitor/app/server/grpc/controller"

	api "github.com/krecu/go-visitor/app/server/grpc/protoc/visitor"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"
)

var (
	Config = viper.New()
	Logger = logrus.New()
)

func init() {

	// enable rand seed
	rand.Seed(time.Now().UnixNano())

	// read & init configure
	Config = viper.New()
	Config.SetEnvPrefix("VISITOR")
	Config.SetConfigName("config")
	Config.AddConfigPath("$HOME/.visitor")

	// check argument config
	CustomConf := ""
	flag.StringVar(&CustomConf, "config", "", "config path")
	flag.Parse()
	if CustomConf != "" {
		fmt.Println(CustomConf)
		Config.AddConfigPath(CustomConf)
	} else {
		Config.AddConfigPath(".")
	}
	Config.AutomaticEnv()
	err := Config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	Logger.Out = bufio.NewWriterSize(os.Stdout, 1024*16)
	Logger.Infof("Конфиг: %s", Config.ConfigFileUsed())
	Logger.Level = logrus.Level(Config.GetInt("app.system.DebugLevel"))

	// добавляем логирование в грейлог
	hook := graylog.NewGraylogHook(
		fmt.Sprintf("%s:%s", Config.GetString("app.graylog.Host"), Config.GetString("app.graylog.Port")),
		map[string]interface{}{"facility": Config.GetString("app.system.Instance")})

	// вырубаем сжатие сообщений ради экономии CPU
	wg := hook.Writer()
	wg.CompressionType = graylog.NoCompress
	wg.CompressionLevel = 0
	hook.SetWriter(wg)
	Logger.AddHook(hook)

	GRPCServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(uint32(Config.GetInt("app.server.MaxConcurrentStreams"))),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: Config.GetDuration("app.server.MaxConnectionIdle") * time.Second,
			Time:              Config.GetDuration("app.server.Time") * time.Second,
			Timeout:           Config.GetDuration("app.server.Timeout") * time.Second,
		}),
	)

	GRPCController, err := controller.New()
	if err != nil {
		Logger.Fatalf("GRPC: %s", err)
	}

	api.RegisterGreeterServer(GRPCServer, GRPCController)
	reflection.Register(GRPCServer)

	// формируем список провайдеров и добавляем в визитор

	// Основная GEO база для RU трафика
	if sp, err := sypexgeo.New(sypexgeo.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
		Weight: 2,
		Name:   "sypexgeo",
	}); err == nil {
		GRPCController.SetGeoProvider(sp)
	} else {
		Logger.Fatalf("SYPEXGEO: %s", err)
	}

	// Основная ГЕО база для не RU трафика
	if mm, err := maxmind.New(maxmind.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/GeoLite2-City.mmdb",
		Weight: 1,
		Name:   "maxmind",
	}); err == nil {
		GRPCController.SetGeoProvider(mm)
	} else {
		Logger.Fatalf("MAXMIND: %s", err)
	}

	// Основная база User-Agent
	if br, err := browscap.New(browscap.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/full_php_browscap.ini",
		Weight: 2,
		Name:   "browscap",
	}); err == nil {
		GRPCController.SetDeviceProvider(br)
	} else {
		Logger.Fatalf("BROWSCAP: %s", err)
	}

	// Словарное определение User-Agent
	if ua, err := uasurfer.New(uasurfer.Option{
		Weight: 1,
		Name:   "uasurfer",
	}); err == nil {
		GRPCController.SetDeviceProvider(ua)
	} else {
		Logger.Fatalf("UASURFER: %s", err)
	}

	if listen, err := net.Listen("tcp", Config.GetString("app.server.listen")); err == nil {
		Logger.Fatalf("GRPC: %s", GRPCServer.Serve(listen))
	} else {
		Logger.Fatalf("GRPC: %s", err)
	}

}
