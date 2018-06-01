package sypexgeo

import (
	"fmt"

	"net"

	"github.com/CossackPyra/pyraconv"
	"github.com/krecu/go-visitor/app/module/provider/geo"
	"github.com/night-codes/go-sypexgeo"
)

type SypexGeo struct {
	conn sypexgeo.SxGEO
	gn   []geo.CountryGeoNames
}

type Option struct {
	Db string
}

func New(opt Option) (proto *SypexGeo, err error) {

	proto = &SypexGeo{
		conn: sypexgeo.New(opt.Db),
	}

	if proto.conn.Version == 0 {
		err = fmt.Errorf("Error load DB: %s", opt.Db)
	}

	proto.gn, err = geo.LoadCountry()

	return
}

func (spx *SypexGeo) Get(ip string) (proto *geo.Model, err error) {

	var (
		res sypexgeo.Result
	)
	ipV4 := net.ParseIP(ip)

	// check correct ip
	if ipV4.To4() != nil && ipV4.IsGlobalUnicast() {
		res, err = spx.conn.Info(ip)
		if err == nil {
			proto = &geo.Model{
				Country: struct {
					Id                 string
					Name               string
					Iso                string
					ISO_3166_1_alpha_3 string
				}{Id: pyraconv.ToString(res.Country.ID), Name: res.Country.NameEn, Iso: res.Country.ISO},
				City: struct {
					Id   string
					Name string
				}{Id: pyraconv.ToString(res.City.ID), Name: res.City.NameEn},
				Region: struct {
					Id   string
					Name string
					Iso  string
				}{Id: pyraconv.ToString(res.Region.ID), Name: res.Region.NameEn, Iso: res.Region.ISO},
				Ip: struct {
					V4 string
					V6 string
				}{V4: ipV4.To4().String(), V6: ipV4.To16().String()},
				Location: struct {
					Latitude  float64
					Longitude float64
					TimeZone  string
				}{Latitude: pyraconv.ToFloat64(res.City.Lat), Longitude: pyraconv.ToFloat64(res.City.Lon), TimeZone: ""},
				Postal: struct{ Code string }{Code: ""},
			}

			if proto.Country.Name != "" {
				for _, c := range spx.gn {
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

func (spx *SypexGeo) Close() {
	spx.Close()
}
