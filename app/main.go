package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gemnasium/logrus-graylog-hook.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	CmdTrigger = make(chan string)
	Config     = viper.New()
	Logger     = logrus.New()
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

	// init logger
	if Config.GetString("app.system.mode") == "prod" {
		Logger.Out = &lumberjack.Logger{
			Filename:   Config.GetString("app.system.LogFile"),
			MaxSize:    500,
			MaxBackups: 3,
			MaxAge:     28,
		}
	} else {
		Logger.Out = os.Stdout
	}

	// добавляем логирование в грейлог
	hook := graylog.NewGraylogHook(
		Config.GetString("app.graylog.host")+":"+Config.GetString("app.graylog.port"),
		map[string]interface{}{"facility": Config.GetString("app.system.instance")})
	Logger.AddHook(hook)

	// init logger level
	Logger.Level = logrus.Level(Config.GetInt("app.system.DebugLevel"))
	Logger.Infof("use config %s", Config.ConfigFileUsed())

	// инициализируем приложение
	Application, err := NewApplication(Config)
	if err != nil {
		panic(err)
	}

	// стартуем приложение
	go Application.Start()

	if Config.GetBool("app.server.grpc.enable") {
		CmdTrigger <- "start-rpc"
	}

	if Config.GetBool("app.server.http.enable") {
		CmdTrigger <- "start-http"
	}

	select {}
}
