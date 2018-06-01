package sypexgeo

import (
	"testing"

	"log"

	"github.com/Pallinder/go-randomdata"
)

func TestSypex_Get(t *testing.T) {
	spx, err := New(Option{
		Db: "./SxGeoMax.dat",
	})

	if err != nil {
		t.Error(err)
	}

	ip := randomdata.IpV4Address()
	p, err := spx.Get(ip)
	if err != nil {
		t.Logf("IP: %s generate err: %s\r\n", ip, err)
	} else {
		log.Printf("%+v", p)
	}
}

func BenchmarkSypex_Get(b *testing.B) {
	spx, err := New(Option{
		Db: "./SxGeoMax.dat",
	})

	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		ip := randomdata.IpV4Address()
		_, err := spx.Get(ip)
		if err != nil {
			b.Logf("IP: %s generate err: %s\r\n", ip, err)
		}
	}
}
