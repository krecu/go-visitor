package wrapper

import (
	"bufio"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/k0kubun/pp"
	"github.com/krecu/go-visitor/app/module/provider/device/uasurfer"
	"github.com/krecu/go-visitor/app/module/provider/geo/maxmind"
	"github.com/krecu/go-visitor/app/module/provider/geo/sypexgeo"
)

func TestWrapper_Parse_Default(t *testing.T) {

	wr := New()

	if sp, err := sypexgeo.New(sypexgeo.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
		Weight: 2,
		Name:   "sypexgeo",
	}); err == nil {
		wr.AddGeoProvider(sp)
	}

	if mm, err := maxmind.New(maxmind.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/GeoLite2-City.mmdb",
		Weight: 1,
		Name:   "maxmind",
	}); err == nil {
		wr.AddGeoProvider(mm)
	}

	if ua, err := uasurfer.New(uasurfer.Option{
		Weight: 1,
		Name:   "uasurfer",
	}); err == nil {
		wr.AddDeviceProvider(ua)
	}

	ua, ip := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36", "79.104.42.249"
	t.Logf("%s, %s", ip, ua)

	info, err := wr.Parse(ip, ua)

	if err != nil {
		t.Errorf("%s, %s: %s", ip, ua, err)
	} else {
		if info == nil {
			t.Errorf("%s, %s: не был распознан", ip, ua)
		}
	}
}

func BenchmarkWrapper_Parse(b *testing.B) {

	wr := New()

	if sp, err := sypexgeo.New(sypexgeo.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
		Weight: 2,
		Name:   "sypexgeo",
	}); err == nil {
		wr.AddGeoProvider(sp)
	}

	if mm, err := maxmind.New(maxmind.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/GeoLite2-City.mmdb",
		Weight: 1,
		Name:   "maxmind",
	}); err == nil {
		wr.AddGeoProvider(mm)
	}

	if ua, err := uasurfer.New(uasurfer.Option{
		Weight: 1,
		Name:   "uasurfer",
	}); err == nil {
		wr.AddDeviceProvider(ua)
	}

	lines := GetTV()
	rand.Seed(time.Now().UnixNano())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ua, ip := lines[rand.Intn(len(lines)-1)], randomdata.IpV4Address()
		d, _ := wr.Parse(ip, ua)
		pp.Println((d.Debug.TimeDevice.Nanoseconds() + d.Debug.TimeGeo.Nanoseconds()))
	}
}

func GetTV() []string {
	file, _ := os.Open("./../../db/tv.ua")
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
