package maxmind

import (
	"testing"

	"log"

	"github.com/Pallinder/go-randomdata"
)

func TestMaxMind_Get(t *testing.T) {
	mm, err := New(Option{
		db: "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/databases/GeoLite2-City.mmdb",
	})

	if err != nil {
		t.Error(err)
	}

	ip := randomdata.IpV4Address()
	p, err := mm.Get(ip)
	if err != nil {
		t.Logf("IP: %s generate err: %s\r\n", ip, err)
	} else {
		log.Printf("%+v", p)
	}
}

func BenchmarkMaxMind_Get(b *testing.B) {
	mm, err := New(Option{
		db: "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/databases/GeoLite2-City.mmdb",
	})

	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		ip := randomdata.IpV4Address()
		_, err := mm.Get(ip)
		if err != nil {
			b.Logf("IP: %s generate err: %s\r\n", ip, err)
		}
	}
}
