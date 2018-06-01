package visitor

import (
	"fmt"
	"net"
	"visitor/model"

	"github.com/CossackPyra/pyraconv"
	"github.com/oschwald/geoip2-golang"
)

var VisitorErrorEmpty = fmt.Errorf("Visitor empty ")

type Visitor struct {
}

func New() (proto *Visitor, err error) {

	proto = &Visitor{}

	return
}

// определение географии
func (v *VisitorService) GetGeo(ip string) (geo *model.Geo, err error) {

	var (
		cacheKey = fmt.Sprintf("geo_%s", ip)
		vSypex   map[string]interface{}
		vMaxMind *geoip2.City
	)

	geo = &model.Geo{}

	if GeoCache, ok := v.app.cache.Get(cacheKey); GeoCache != nil && ok {
		if value, ok := GeoCache.(*model.Geo); ok {
			geo = value
			return
		}
	}

	// считываем ГЕО данные по SyPex базе
	if vSypex, err = v.app.sypexgeo.GetCityFull(ip); err == nil {
		if country, ok := vSypex["country"].(map[string]interface{}); ok {
			geo.Country.Id = uint(pyraconv.ToInt64(country["id"]))
			geo.Country.Name = pyraconv.ToString(country["name_en"])
			geo.Country.NameRu = pyraconv.ToString(country["name_ru"])
			geo.Country.Iso = pyraconv.ToString(country["iso"])
			geo.Country.Iso3166_1_alpha_3 = model.ISO_3166_1_alpha_3Mapping(pyraconv.ToString(country["iso"]))
			geo.Country.Mapping = model.CountryMapping(geo.Country.Iso)
			geo.Location.TimeZone = pyraconv.ToString(country["timezone"])
		}
		if region, ok := vSypex["region"].(map[string]interface{}); ok {
			geo.Region.Id = uint(pyraconv.ToInt64(region["id"]))
			geo.Region.Name = pyraconv.ToString(region["name_en"])
			geo.Region.NameRu = pyraconv.ToString(region["name_ru"])
			geo.Region.Iso = pyraconv.ToString(region["iso"])
			geo.Region.Mapping = model.RegionMapping(geo.Region.NameRu)
		}
		if city, ok := vSypex["city"].(map[string]interface{}); ok {
			geo.City.Id = uint(pyraconv.ToInt64(city["id"]))
			geo.City.Name = pyraconv.ToString(city["name_en"])
			geo.City.NameRu = pyraconv.ToString(city["name_ru"])
			geo.Location.Latitude = pyraconv.ToFloat32(city["lat"])
			geo.Location.Longitude = pyraconv.ToFloat32(city["lon"])
			geo.City.Mapping = model.CityMapping(geo.City.NameRu)
		}
	}

	// дополянем и перероверяем ГЕО по MaxMind
	if vMaxMind, err = v.app.maxmind.City(net.ParseIP(ip)); err == nil {

		if country := vMaxMind.Country; country.GeoNameID != 0 && geo.Country.Id == 0 {
			geo.Country.Id = uint(pyraconv.ToInt64(country.GeoNameID))
			geo.Country.Name = pyraconv.ToString(country.Names["en"])
			geo.Country.NameRu = pyraconv.ToString(country.Names["ru"])
			geo.Country.Iso = pyraconv.ToString(country.IsoCode)
			geo.Country.Mapping = model.CountryMapping(country.Names["ru"])
		}
		if city := vMaxMind.City; city.GeoNameID != 0 && geo.City.Id == 0 {
			geo.City.Id = uint(pyraconv.ToInt64(city.GeoNameID))
			geo.City.Name = pyraconv.ToString(city.Names["en"])
			geo.City.NameRu = pyraconv.ToString(city.Names["ru"])
			geo.City.Mapping = model.CityMapping(city.Names["ru"])
		}
		if location := vMaxMind.Location; location.TimeZone != "" && geo.Location.TimeZone == "" {
			geo.Location.Latitude = pyraconv.ToFloat32(location.Latitude)
			geo.Location.Longitude = pyraconv.ToFloat32(location.Longitude)
			geo.Location.TimeZone = pyraconv.ToString(location.TimeZone)
		}
		if postal := vMaxMind.Postal; postal.Code != "" {
			geo.Postal.Code = pyraconv.ToString(postal.Code)
		}
	}

	if err == nil && geo != nil {
		v.app.cache.SetDefault(cacheKey, geo)

	}

	return
}

func (v *Visitor) Identification(id string, ip string, ua string, extra map[string]interface{}) (proto *model.Visitor, err error) {

	return
}
