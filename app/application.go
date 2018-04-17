package main

import (
	"strings"

	"time"

	"github.com/CossackPyra/pyraconv"
	"github.com/aerospike/aerospike-client-go"
	"github.com/krecu/browscap_go"
	//sypexgeo "github.com/krecu/go-sypexgeo"
	sypexgeo "github.com/night-codes/go-sypexgeo"
	"github.com/oschwald/geoip2-golang"
	cache "github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/ua-parser/uap-go/uaparser"
)

type App struct {
	config    *viper.Viper
	rpc       *RpcService
	http      *HttpService
	visitor   *VisitorService
	aerospike *aerospike.Client
	cache     *cache.Cache
	sypexgeo  sypexgeo.SxGEO
	maxmind   *geoip2.Reader
	uaparse   *uaparser.Parser
	browscap  *browscap_go.Browser
}

func NewApplication(config *viper.Viper) (proto *App, err error) {

	proto = &App{
		config: config,
	}

	hosts := []*aerospike.Host{}

	for _, h := range config.GetStringSlice("app.aerospike.Hosts") {
		arg := strings.Split(h, ":")
		hosts = append(hosts, &aerospike.Host{
			Name: pyraconv.ToString(arg[0]),
			Port: int(pyraconv.ToInt64(arg[1])),
		})
	}

	proto.aerospike, err = aerospike.NewClientWithPolicyAndHost(aerospike.NewClientPolicy(), hosts...)
	if err != nil {
		return
	}

	proto.cache = cache.New(
		config.GetDuration("app.cache.DefaultExpiration")*time.Minute,
		config.GetDuration("app.cache.CleanupInterval")*time.Minute,
	)

	proto.sypexgeo = sypexgeo.New(config.GetString("app.database.SxGeoCity"))

	proto.maxmind, err = geoip2.Open(config.GetString("app.database.MaxMind"))
	if err != nil {
		return
	}

	if err = browscap_go.InitBrowsCap(
		config.GetString("app.database.BrowsCap"),
		true,
		config.GetDuration("app.cache.DefaultExpiration")*time.Hour,
		config.GetDuration("app.cache.CleanupInterval")*time.Hour,
	); err != nil {
		return
	}

	// добавляем в контейнер клиента к очереди
	proto.visitor, err = NewVisitorService(proto)
	if err != nil {
		return
	}

	// добавляем в контейнер веб сервер
	proto.http, err = NewHttpService(proto)
	if err != nil {
		return
	}

	// добавляем в контейнер веб сервер
	proto.rpc, err = NewRpcService(proto)
	if err != nil {
		return
	}

	return
}

func (a *App) Config() *viper.Viper {
	return a.config
}

func (a *App) Start() {
	for {
		select {
		case cmd := <-CmdTrigger:
			Logger.Infof("Выолняем команду: %s", cmd)

			switch cmd {
			case "start-rpc":
				a.rpc.Start()
				break
			case "start-http":
				a.http.Start()
				break
			}
		}
	}
}
