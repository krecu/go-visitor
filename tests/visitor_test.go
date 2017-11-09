package visitor_test

import (
	"testing"
	"github.com/krecu/go-visitor"
	"os"
	"bufio"
	"math/rand"
	"github.com/Pallinder/go-randomdata"
	"time"
	"fmt"
)

const (
	UA = "Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30"
	IP = "79.104.42.249"
)

func TestVisitor_GetCity(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	info, err := v.GetCity(IP)

	if err != nil || info.Name == "" {
		t.Fatal(info)
	} else {
		t.Log(info)
	}
}

func TestVisitor_GetCountry(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	info, err := v.GetCountry(IP)

	if err != nil || info.Name == "" {
		t.Fatal(info)
	} else {
		t.Log(info)
	}
}

func TestVisitor_GetRegion(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	info, err := v.GetRegion(IP)

	if err != nil || info.Name == "" {
		t.Fatal(info)
	} else {
		t.Log(info)
	}
}

func TestVisitor_GetLocation(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	info, err := v.GetLocation(IP)

	if err != nil {
		t.Fatal(info)
	} else {
		t.Log(info)
	}
}

func TestVisitor_GetPostal(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	info, err := v.GetPostal(IP)

	if err != nil {
		t.Fatal(info)
	} else {
		t.Log(info)
	}
}

func TestVisitor_GetBrowser(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	info, err := v.GetBrowser(UA)

	if err != nil || info.Name == ""  {
		t.Fatal(info)
	} else {
		t.Log(info)
	}
}

func TestVisitor_GetDevice(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	info, err := v.GetDevice(UA)

	if err != nil || info.Name == "" {
		t.Fatal(info)
	} else {
		t.Log(info)
	}
}

func TestVisitor_Identify(t *testing.T) {

	v, _ := visitor.New()
	v.Debug = true

	file, err := os.Open("./fixtures/user-agent.txt"); if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(420)

	total := time.Now()

	info, err := v.Identify(randomdata.IpV4Address(), lines[rand.Intn(len(lines))])
	fmt.Println("Identify: Execute time " + time.Since(total).String())

	if err != nil {
		t.Fatal(info)
	}
}


func BenchmarkVisitor_Identify(b *testing.B) {
	v, _ := visitor.New()

	file, err := os.Open("./fixtures/user-agent.txt"); if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(42)

	for n := 0; n < b.N; n++ {
		ua := lines[rand.Intn(len(lines))]
		ip := randomdata.IpV4Address()
		v.Identify(ip, ua)
	}
}