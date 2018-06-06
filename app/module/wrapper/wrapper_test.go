package wrapper

import (
	"bufio"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/k0kubun/pp"

	providerDevice "github.com/krecu/go-visitor/app/module/provider/device"
	"github.com/krecu/go-visitor/app/module/provider/device/browscap"
	"github.com/krecu/go-visitor/app/module/provider/device/uasurfer"
	providerGeo "github.com/krecu/go-visitor/app/module/provider/geo"
	"github.com/krecu/go-visitor/app/module/provider/geo/maxmind"
	"github.com/krecu/go-visitor/app/module/provider/geo/sypexgeo"
)

func TestWrapper_Parse(t *testing.T) {
	var (
		geo    []providerGeo.Geo
		device []providerDevice.Device
	)

	if sp, err := sypexgeo.New(sypexgeo.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
		Weight: 2,
		Name:   "sypexgeo",
	}); err == nil {
		geo = append(geo, sp)
	}

	if mm, err := maxmind.New(maxmind.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/GeoLite2-City.mmdb",
		Weight: 1,
		Name:   "maxmind",
	}); err == nil {
		geo = append(geo, mm)
	}

	if br, err := browscap.New(browscap.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/full_php_browscap.ini",
		Weight: 2,
		Name:   "browscap",
	}); err == nil {
		device = append(device, br)
	}

	if ua, err := uasurfer.New(uasurfer.Option{
		Weight: 1,
		Name:   "uasurfer",
	}); err == nil {
		device = append(device, ua)
	}

	wr, err := New(geo, device)
	if err != nil {
		t.Errorf("%s", err)
	}

	//ua, ip := GetUaIp()
	ua, ip := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36", "79.104.42.249"
	//ua, ip := "AppleCoreMedia/1.0.0.12B466 (Apple TV; U; CPU OS 8_1_3 like Mac OS X; en_us)", "79.104.42.249"
	t.Logf("%s, %s", ip, ua)

	info, err := wr.Parse(ip, ua)
	if err != nil {
		t.Errorf("%s, %s: %s", ip, ua, err)
	} else {
		pp.Println(info)
		pp.Println(info.Debug.TimeGeo.Seconds())
		pp.Println(info.Debug.TimeDevice.Seconds())
	}
}

func BenchmarkWrapper_Parse(b *testing.B) {

	var (
		geo    []providerGeo.Geo
		device []providerDevice.Device
	)

	if sp, err := sypexgeo.New(sypexgeo.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
		Weight: 2,
		Name:   "sypexgeo",
	}); err == nil {
		geo = append(geo, sp)
	}

	if mm, err := maxmind.New(maxmind.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/GeoLite2-City.mmdb",
		Weight: 1,
		Name:   "maxmind",
	}); err == nil {
		geo = append(geo, mm)
	}

	if br, err := browscap.New(browscap.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/full_php_browscap.ini",
		Weight: 2,
		Name:   "browscap",
	}); err == nil {
		device = append(device, br)
	}

	if ua, err := uasurfer.New(uasurfer.Option{
		Weight: 1,
		Name:   "uasurfer",
	}); err == nil {
		device = append(device, ua)
	}

	wr, err := New(geo, device)
	if err != nil {
		b.Errorf("%s", err)
	}

	//ua, ip := GetUaIp()
	ua, ip := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36", "79.104.42.249"
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
