package server

import (
	"os"
	"encoding/json"
	"log"
	"time"
)

type ConfGrayLog struct {
	Host string
	Port int
}

type ConfAeroSpike struct {
	Host string
	Port int
	Ns string
	Db string
	Timeout time.Duration
	Ttl uint32
}

type Config struct {
	Cpu    int
	Listener string
	Log 	string
	Debug 	bool
	Buffer  int
	Db 	string
	Logger ConfGrayLog
	AeroSpike ConfAeroSpike
}

// загрузка конфига
func NewConfig(path string) (*Config){

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening conf file: %v", err)
	}

	conf := new(Config)
	err = json.NewDecoder(file).Decode(&conf)

	if err != nil {
		log.Fatalf("error parsing conf file: %v", err)
	}
	return conf
}