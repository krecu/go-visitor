package wrapper

import (
	"bufio"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
)

func TestWrapper_Parse(t *testing.T) {

	wr, err := New()
	if err != nil {
		t.Errorf("%s", err)
	}

	ua, ip := GetUaIp()
	t.Logf("%s, %s", ip, ua)

	_, err = wr.Parse(ip, ua)
	if err != nil {
		t.Errorf("%s, %s: %s", ip, ua, err)
	}
}

func BenchmarkWrapper_Parse(b *testing.B) {

	wr, err := New()
	if err != nil {
		b.Errorf("%s", err)
	}

	ua, ip := GetUaIp()
	b.Logf("%s, %s", ip, ua)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = wr.Parse(ip, ua)
		if err != nil {
			b.Errorf("%s, %s: %s", ip, ua, err)
		}
	}
}

func GetUaIp() (ua string, ip string) {
	file, _ := os.Open("./user-agent.txt")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())
	ua = lines[rand.Intn(len(lines)-1)]
	ip = randomdata.IpV4Address()

	return
}
