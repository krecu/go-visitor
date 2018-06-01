package wrapper

import (
	providerDevice "github.com/krecu/go-visitor/app/module/provider/device"
	"github.com/krecu/go-visitor/app/module/provider/device/browscap"
	providerGeo "github.com/krecu/go-visitor/app/module/provider/geo"
	"github.com/krecu/go-visitor/app/module/provider/geo/sypexgeo"
)

type Wrapper struct {
	geo    providerGeo.Geo
	device providerDevice.Device
}

func New() (proto *Wrapper, err error) {

	proto = &Wrapper{}

	proto.geo, err = sypexgeo.New(sypexgeo.Option{
		Db: "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
	})

	proto.device, err = browscap.New(browscap.Option{
		Db: "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/full_php_browscap.ini",
	})

	return
}

func (wr *Wrapper) Parse(ip string, ua string) (proto *Model, err error) {

	_, err = wr.geo.Get(ip)
	if err != nil {
		return
	}

	_, err = wr.device.Get(ua)
	if err != nil {
		return
	}

	return
}
