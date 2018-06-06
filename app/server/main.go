package main

import (
	"math/rand"
	"runtime"

	"time"

	"fmt"

	"flag"

	"github.com/krecu/go-visitor/app/module/provider/device/browscap"
	"github.com/krecu/go-visitor/app/module/provider/device/uasurfer"
	"github.com/krecu/go-visitor/app/module/provider/geo/maxmind"
	"github.com/krecu/go-visitor/app/module/provider/geo/sypexgeo"
	"github.com/krecu/go-visitor/app/server/grpc"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	//Logger.Out = bufio.NewWriterSize(os.Stdout, 1024*16)
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

	RpcServer, err := grpc.New(grpc.Options{
		Listen:               Config.GetString("app.server.listen"),
		Timeout:              Config.GetDuration("app.server.Timeout"),
		Time:                 Config.GetDuration("app.server.Time"),
		MaxConnectionIdle:    Config.GetDuration("app.server.MaxConnectionIdle"),
		MaxConcurrentStreams: uint32(Config.GetInt("app.server.MaxConcurrentStreams")),
		Db: struct {
			Hosts             []string
			ConnectionTimeout time.Duration
			GetTimeout        time.Duration
			WriteTimeout      time.Duration
			Expiration        uint32
			Set               string
			NameSpace         string
		}{
			Hosts:             Config.GetStringSlice("app.aerospike.Hosts"),
			ConnectionTimeout: Config.GetDuration("app.aerospike.ConnectionTimeout"),
			GetTimeout:        Config.GetDuration("app.aerospike.GetTimeout"),
			WriteTimeout:      Config.GetDuration("app.aerospike.WriteTimeout"),
			Expiration:        uint32(Config.GetInt("app.aerospike.Expiration")),
			Set:               Config.GetString("app.aerospike.Set"),
			NameSpace:         Config.GetString("app.aerospike.NameSpace"),
		},
		Visitor: struct {
			SypexGeo sypexgeo.Option
			MaxMind  maxmind.Option
			BrowsCap browscap.Option
			UaSurfer uasurfer.Option
		}{
			SypexGeo: sypexgeo.Option{
				Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
				Weight: 2,
				Name:   "sypexgeo",
			},
			MaxMind: maxmind.Option{
				Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/GeoLite2-City.mmdb",
				Weight: 1,
				Name:   "maxmind",
			},
			BrowsCap: browscap.Option{
				Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/full_php_browscap.ini",
				Weight: 2,
				Name:   "browscap",
			},
			UaSurfer: uasurfer.Option{
				Weight: 1,
				Name:   "uasurfer",
			},
		},
	})

	if err != nil {
		Logger.Fatalf("GRPC: %s", err)
	}

	Logger.Fatal(RpcServer.Start())

}
