package visitor

import (
	"bufio"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
)

func TestVisitor_Identification(t *testing.T) {

	var (
		err error
		v   *Visitor
	)

	v, err = New()
	if err != nil {
		t.Fail()
	}

	file, _ := os.Open("./user-agent.txt")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	rand.Seed(time.Now().UnixNano())

	ua := lines[rand.Intn(len(lines)-1)]
	ip := randomdata.IpV4Address()

	_, err = v.Identification("test_user", ip, ua, nil)
	if err != nil {
		t.Fail()
	}
}
