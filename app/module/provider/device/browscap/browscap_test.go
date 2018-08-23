package browscap

import (
	"bufio"
	"os"
	"testing"

	"github.com/k0kubun/pp"
)

func TestBrowsCap_Get_TV(t *testing.T) {
	bc, err := New(Option{
		Db: "./../../../../db/full_php_browscap_1.ini",
	})

	if err != nil {
		t.Error(err)
	}

	list := GetUaIp()

	for _, s := range list {
		d, _ := bc.Get(s)
		//if d.Device.Type
		// != "TV Device" {
		pp.Println(d.Device, s)
		//}
	}
}

func BenchmarkBrowsCap_Get(b *testing.B) {
	bc, err := New(Option{
		Db: "./full_php_browscap.ini",
	})

	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		bc.Get("Mozilla/5.0 (iPhone; CPU iPhone OS 11_2_6 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) Mobile/15D100 Instagram 28.0.0.12.285 (iPhone9,3; iOS 11_2_6; ru_RU; ru-RU; scale=2.00; gamut=wide; 750x1334)")
	}
}

func GetUaIp() []string {
	file, _ := os.Open("./tv.ua")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
