package main

import (
	"flag"
	"runtime"
	"log"
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/krecu/go-visitor/server"
	"github.com/krecu/go-visitor/cache"
	"github.com/krecu/go-visitor"
)

var (
	Conf *server.Config
	Core *server.Core
)


// Стартуем приложение
func main() {

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

	log.Fatal(NewHttp(Conf.Listener))
}