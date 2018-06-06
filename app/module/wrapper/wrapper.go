package wrapper

import (
	"time"

	"sort"

	providerDevice "github.com/krecu/go-visitor/app/module/provider/device"
	providerGeo "github.com/krecu/go-visitor/app/module/provider/geo"
)

type Wrapper struct {
	geo    []providerGeo.Geo
	device []providerDevice.Device
}

func New() (proto *Wrapper) {
	proto = &Wrapper{}
	return
}

func (wr *Wrapper) AddGeoProvider(geo providerGeo.Geo) {

	for _, t := range wr.geo {
		if t.Name() == geo.Name() {
			return
		}
	}

	wr.geo = append(wr.geo, geo)
}

func (wr *Wrapper) AddDeviceProvider(device providerDevice.Device) {

	for _, t := range wr.geo {
		if t.Name() == device.Name() {
			return
		}
	}

	wr.device = append(wr.device, device)
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

	if len(wr.geo) > 0 {
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
	}

	if len(wr.device) > 0 {
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
	}
	return
}
