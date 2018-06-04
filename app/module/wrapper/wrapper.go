package wrapper

import (
	"time"

	"sort"

	providerDevice "github.com/krecu/go-visitor/app/module/provider/device"
	"github.com/krecu/go-visitor/app/module/provider/device/browscap"
	"github.com/krecu/go-visitor/app/module/provider/device/uasurfer"
	providerGeo "github.com/krecu/go-visitor/app/module/provider/geo"
	"github.com/krecu/go-visitor/app/module/provider/geo/maxmind"
	"github.com/krecu/go-visitor/app/module/provider/geo/sypexgeo"
)

type Wrapper struct {
	geo    []providerGeo.Geo
	device []providerDevice.Device
}

func New() (proto *Wrapper, err error) {

	proto = &Wrapper{}

	if sp, err := sypexgeo.New(sypexgeo.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/SxGeoMax.dat",
		Weight: 2,
		Name:   "sypexgeo",
	}); err == nil {
		proto.geo = append(proto.geo, sp)
	}

	if mm, err := maxmind.New(maxmind.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/GeoLite2-City.mmdb",
		Weight: 1,
		Name:   "maxmind",
	}); err == nil {
		proto.geo = append(proto.geo, mm)
	}

	if br, err := browscap.New(browscap.Option{
		Db:     "/Users/kretsu/Work/Go/src/github.com/krecu/go-visitor/app/db/full_php_browscap.ini",
		Weight: 2,
		Name:   "browscap",
	}); err == nil {
		proto.device = append(proto.device, br)
	}

	if ua, err := uasurfer.New(uasurfer.Option{
		Weight: 1,
		Name:   "uasurfer",
	}); err == nil {
		proto.device = append(proto.device, ua)
	}

	return
}

func (wr *Wrapper) Parse(ip string, ua string) (proto *Model, err error) {

	var (
		infoGeo    *providerGeo.Model
		infoDevice *providerDevice.Model
	)

	proto = &Model{
		Personal: Personal{
			Ua: ua,
		},
	}

	_InfotGeo := time.Now()

	sort.Sort(providerGeo.OrderProvider(wr.geo))
	// перебираем гео провайдеров до успеха
	for _, p := range wr.geo {
		infoGeo, err = p.Get(ip)
		if err != nil {
			continue
		} else {
			proto.Debug.ProviderGeo = p.Name()
			break
		}
	}
	proto.Debug.TimeGeo = time.Since(_InfotGeo)

	if err != nil {
		return
	} else {
		proto.Country = Country{
			Id:                infoGeo.Country.Id,
			Name:              infoGeo.Country.Name,
			Iso:               infoGeo.Country.Iso,
			Iso3166_1_alpha_3: infoGeo.Country.ISO_3166_1_alpha_3,
		}

		proto.City = City{
			Id:   infoGeo.City.Id,
			Name: infoGeo.City.Name,
		}

		proto.Region = Region{
			Id:   infoGeo.Region.Id,
			Name: infoGeo.Region.Name,
			Iso:  infoGeo.Region.Iso,
		}

		proto.Location = Location{
			TimeZone:  infoGeo.Location.TimeZone,
			Longitude: infoGeo.Location.Longitude,
			Latitude:  infoGeo.Location.Latitude,
		}

		proto.Postal = Postal{
			Code: infoGeo.Postal.Code,
		}

		proto.Ip = Ip{
			V4: infoGeo.Ip.V4,
			V6: infoGeo.Ip.V6,
		}
	}

	_InfoDevice := time.Now()

	sort.Sort(providerDevice.OrderProvider(wr.device))

	// перебираем гео провайдеров до успеха
	for _, p := range wr.device {
		infoDevice, err = p.Get(ua)
		if err != nil {
			continue
		} else {
			proto.Debug.ProviderDevice = p.Name()
			break
		}
	}

	proto.Debug.TimeDevice = time.Since(_InfoDevice)

	if err != nil {
		return
	} else {
		proto.Device = Device{
			Name:  infoDevice.Device.Name,
			Brand: infoDevice.Device.Brand,
			Type:  infoDevice.Device.Type,
		}

		proto.Browser = Browser{
			Name:    infoDevice.Browser.Name,
			Type:    infoDevice.Browser.Type,
			Version: infoDevice.Browser.Version,
		}

		proto.Platform = Platform{
			Name:    infoDevice.Platform.Name,
			Short:   infoDevice.Platform.Short,
			Version: infoDevice.Platform.Version,
		}
	}
	return
}
