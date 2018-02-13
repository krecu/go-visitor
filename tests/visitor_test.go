package tests

import (
	"testing"
	"os"
	"bufio"
	"math/rand"
	"github.com/Pallinder/go-randomdata"
	"time"
	"fmt"
	"github.com/krecu/go-visitor/client/rpc-client"
)

const (
	UA = "Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30"
	IP = "79.104.42.249"
)


func TestVisitor_Identify(t *testing.T) {

	v, _ := visitor_rpc_client.New([]string{"127.0.0.1:8081"})

	file, err := os.Open("./fixtures/user-agent.txt"); if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())

	total := time.Now()
	info, err := v.Post("test_user", randomdata.IpV4Address(), lines[rand.Intn(len(lines))], nil)
	fmt.Println("Identify: Execute time " + time.Since(total).String())

	if err != nil {
		t.Fatal(info)
	}
}


func BenchmarkVisitor_Identify(b *testing.B) {
	v, _ := visitor_rpc_client.New([]string{"127.0.0.1:8081"})

	file, err := os.Open("./fixtures/user-agent.txt"); if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())

	for n := 0; n < b.N; n++ {
		ua := lines[rand.Intn(len(lines))]
		ip := randomdata.IpV4Address()
		v.Post("test_user", ip, ua, nil)
	}
}