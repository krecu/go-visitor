package maxmind

import (
	"fmt"

	"net"

	"github.com/CossackPyra/pyraconv"
	"github.com/krecu/go-visitor/app/module/provider/geo"
	"github.com/oschwald/geoip2-golang"
)

type MaxMind struct {
	conn *geoip2.Reader
	gn   []geo.CountryGeoNames
}

type Option struct {
	db string
}

func New(opt Option) (proto *MaxMind, err error) {

	proto = &MaxMind{}
	proto.conn, err = geoip2.Open(opt.db)
	if err != nil {
		return
	}
	proto.gn, err = geo.LoadCountry()

	return
}

func (mm *MaxMind) Get(ip string) (proto *geo.Model, err error) {

	var (
		res *geoip2.City
	)
	ipV4 := net.ParseIP(ip)

	// check correct ip
	if ipV4.To4() != nil && ipV4.IsGlobalUnicast() {
		res, err = mm.conn.City(ipV4)
		if err == nil {
			proto = &geo.Model{
				Country: struct {
					Id                 string
					Name               string
					Iso                string
					ISO_3166_1_alpha_3 string
				}{Id: pyraconv.ToString(res.Country.GeoNameID), Name: res.Country.Names["en"], Iso: res.Country.IsoCode},
				City: struct {
					Id   string
					Name string
				}{Id: pyraconv.ToString(res.City.GeoNameID), Name: res.City.Names["en"]},
				Region: struct {
					Id   string
					Name string
					Iso  string
				}{Id: "", Name: "", Iso: ""},
				Ip: struct {
					V4 string
					V6 string
				}{V4: ipV4.To4().String(), V6: ipV4.To16().String()},
				Location: struct {
					Latitude  float64
					Longitude float64
					TimeZone  string
				}{Latitude: pyraconv.ToFloat64(res.Location.Latitude), Longitude: pyraconv.ToFloat64(res.Location.Longitude), TimeZone: res.Location.TimeZone},
				Postal: struct{ Code string }{Code: res.Postal.Code},
			}

			if proto.Country.Name != "" {
				for _, c := range mm.gn {
					if c.Name == proto.Country.Name {
						proto.Country.Id = c.Id
						proto.Country.ISO_3166_1_alpha_3 = c.Iso3
					}
				}
			}
		}
	} else {
		err = fmt.Errorf("IP is a not global unicast address.")
	}

	return
}

func (mm *MaxMind) Close() {
	mm.Close()
}
