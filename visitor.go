package main

import (

	"os"
	"bufio"
	"math/rand"
	"github.com/Pallinder/go-randomdata"
	"time"
	"github.com/krecu/go-visitor/client/rpc-client"
	"github.com/paulbellamy/ratecounter"
	"log"
)

const (
	UA = "Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30"
	IP = "79.104.42.249"
)


func main() {


	file, _ := os.Open("/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/tests/fixtures/user-agent.txt")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())
	counter_rmq := ratecounter.NewRateCounter(1 * time.Second)

	// включаем счетчики состояния
	go func() {
		for range time.Tick(1 * time.Second) {
			log.Println(counter_rmq.Rate())
		}
	}()

	go func() {
		for {
			ua := lines[rand.Intn(len(lines))]
			ip := randomdata.IpV4Address()
			v, _ := visitor_rpc_client.New([]string{"127.0.0.1:8081"})
			v.Post("test_user", ip, ua, nil)
			counter_rmq.Incr(1)
		}
	}()

	select {}
}