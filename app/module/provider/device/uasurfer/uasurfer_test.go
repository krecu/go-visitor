package uasurfer

import (
	"bufio"
	"os"
	"testing"

	"math/rand"
	"time"

	"github.com/k0kubun/pp"
)

func TestUaSurfer_Get_TV(t *testing.T) {
	uas, err := New(Option{})

	if err != nil {
		t.Error(err)
	}

	list := GetTV()

	for _, s := range list {
		d, _ := uas.Get(s)
		pp.Println(d.Device.Type, d.Platform.Name, d.Browser.Name)
	}
}

func TestUaSurfer_Get_Tablet(t *testing.T) {
	uas, err := New(Option{})

	if err != nil {
		t.Error(err)
	}

	list := GetTablet()

	for _, s := range list {
		d, _ := uas.Get(s)
		pp.Println(d.Device.Type, d.Platform.Name, d.Browser.Name)
	}
}

func BenchmarkUap_Get(b *testing.B) {
	uas, err := New(Option{})

	if err != nil {
		b.Error(err)
	}

	lines := GetTV()
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		ua := lines[rand.Intn(len(lines)-1)]
		uas.Get(ua)
	}
}

func GetTV() []string {
	file, _ := os.Open("./../../../../db/tv.ua")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func GetTablet() []string {
	file, _ := os.Open("./../../../../db/tablet.ua")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
