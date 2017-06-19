package visitor

import (
	"github.com/digitalcrab/browscap_go"
	"github.com/krecu/go-visitor/model"
	"errors"
	"time"
	"fmt"
	"github.com/avct/uasurfer"
	"github.com/oschwald/geoip2-golang"
	"github.com/mirrr/go-sypexgeo"
	"net"
	"github.com/CossackPyra/pyraconv"
)


type Visitor struct {
	Path string
	Cache map[string]interface{}
	Debug bool
	Buffer int
	Processor map[string]interface{}
}

func New(debug bool, path string, buffer int) (info *Visitor, err error) {

	dbPath := path + "/db"

	// init browscap
	if browscap_go.InitializedVersion() == "" {
		if err = browscap_go.InitBrowsCap(dbPath + "/browscap.ini", false); err != nil {
			panic(err)
		}
	}

	processors := make(map[string]interface{})

	// init geo processor
	processors["sypexgeo"] = sypexgeo.New(dbPath+ "/SxGeoCity.dat")
	processors["maxmind"], err = geoip2.Open(dbPath + "/GeoLite2-City.mmdb"); if err != nil {
		panic(err)
	}

	info = &Visitor{
		Debug: debug,
		Path: path,
		Buffer: buffer,
		Processor: processors,
		Cache: make(map[string]interface{}),
	}

	return
}

/**
 Clear cache and close processors
 */
func (v *Visitor) Close() {
	v.Processor["maxmind"].(*geoip2.Reader).Close()
	v.Cache = make(map[string]interface{})
}

/**
 Add to cache
 */
func (v *Visitor) SetCache(key string, val interface{}) {
	if len(v.Cache) > v.Buffer {
		v.Cache = make(map[string]interface{})
		v.Cache[key] = val
	} else {
		v.Cache[key] = val
	}
}

/**
 Get from cache
 */
func (v *Visitor) GetCache(key string) (val interface{}) {
	val, ok := v.Cache[key]; if !ok {
		val = nil
	}
	return
}

/**
 Get aggregation idetify user information
 */
func (v *Visitor) Identify (ip string, ua string) (visitor model.Visitor, err error) {

	total := time.Now()

	visitor = model.Visitor{}

	visitor.Platform, _ = v.GetPlatform(ua)
	visitor.Device, _ = v.GetDevice(ua)
	visitor.Browser, _ = v.GetBrowser(ua)
	visitor.Personal, _ = v.GetPersonal(ua)

	if net.ParseIP(ip) != nil {
		visitor.Ip, _ = v.GetIp(ip)
		visitor.Country, _ = v.GetCountry(ip)
		visitor.City, _ = v.GetCity(ip)
		visitor.Region, _ = v.GetRegion(ip)
		visitor.Location, _ = v.GetLocation(ip)
		visitor.Postal, _ = v.GetPostal(ip)
	}

	if v.Debug {
		fmt.Println("Identify: Execute time " + time.Since(total).String())
	}

	return
}

/**
 Get SyPexGeo info by ip
 */
func (v *Visitor) GetSyPexGeo(ip string) (info map[string]interface{}, err error) {

	total := time.Now()

	// All query by user agent cached
	info, ok := v.GetCache(ip + "_sypexgeo").(map[string]interface{}); if !ok {

		info, err = v.Processor["sypexgeo"].(*sypexgeo.SxGEO).GetCityFull(ip)

		if err != nil || info == nil {
			err = errors.New("IP not found " + ip)
			return
		} else {
			v.SetCache(ip + "_sypexgeo", info)
		}
	}

	if v.Debug {
		fmt.Println("GetSyPexGeo: Execute time " + time.Since(total).String())
	}
	return
}

/**
 Get BrowsCap info by user agent
 */
func (v *Visitor) GetMaxMind(ip string) (info *geoip2.City, err error) {

	total := time.Now()

	// All query by user agent cached
	info, ok := v.GetCache(ip + "_maxmind").(*geoip2.City); if !ok {
		info, err = v.Processor["maxmind"].(*geoip2.Reader).City(net.ParseIP(ip))

		if err != nil || info == nil {
			return
		} else {
			v.SetCache(ip + "_maxmind", info)
		}
	}

	if v.Debug {
		fmt.Println("GetMaxMind: Execute time " + time.Since(total).String())
	}
	return
}

/**
 Get BrowsCap info by user agent
 */
func (v *Visitor) GetBrowsCap(ua string) (browscap *browscap_go.Browser, err error) {

	total := time.Now()

	// All query by user agent cached
	browscap, ok := v.GetCache(ua + "_browscap").(*browscap_go.Browser); if !ok {

		browscap, ok = browscap_go.GetBrowser(ua)

		if !ok || browscap == nil {
			err = errors.New("Browser not found " + ua)
			return
		}

		v.SetCache(ua + "_browscap", browscap)
	}

	if v.Debug {
		fmt.Println("GetBrowsCap: Execute time " + time.Since(total).String())
	}
	return
}

/**
 Get City info by ip
 */
func (v *Visitor) GetCity(ip string) (city model.City, err error) {


	mmRecord, err := v.GetMaxMind(ip); if err == nil {
		city.Name = mmRecord.City.Names["en"]
		city.NameRu = mmRecord.City.Names["ru"]
		city.Id = mmRecord.City.GeoNameID
	}

	if err != nil || city.Name == "" || mmRecord.Country.IsoCode == "RU" || mmRecord.Country.IsoCode == "UA" {

		spRecord, err := v.GetSyPexGeo(ip); if err == nil {
			raw := spRecord["city"].(map[string]interface{})
			city.Name = raw["name_en"].(string)
			city.NameRu = raw["name_ru"].(string)
			city.Id = uint(pyraconv.ToInt64(raw["id"]))
		}
	}

	return
}

/**
 Get Country info by ip
 */
func (v *Visitor) GetCountry(ip string) (country model.Country, err error) {

	mmRecord, err := v.GetMaxMind(ip); if err == nil {
		country.Name = mmRecord.Country.Names["en"]
		country.NameRu = mmRecord.Country.Names["ru"]
		country.Id = mmRecord.Country.GeoNameID
		country.Iso = mmRecord.Country.IsoCode
	}

	if err !=nil || mmRecord.Country.IsoCode == "RU" || mmRecord.Country.IsoCode == "UA" {

		spRecord, err := v.GetSyPexGeo(ip); if err == nil {
			raw := spRecord["country"].(map[string]interface{})
			country.Name = raw["name_en"].(string)
			country.NameRu = raw["name_ru"].(string)
			country.Id = uint(pyraconv.ToInt64(raw["id"]))
			country.Iso = raw["iso"].(string)
		}
	}

	return
}

/**
 Get Region info by ip
 */
func (v *Visitor) GetRegion(ip string) (region model.Region, err error) {

	spRecord, err := v.GetSyPexGeo(ip); if err == nil {
		raw := spRecord["region"].(map[string]interface{})
		if raw["name_en"] != nil {
			region.Name = raw["name_en"].(string)
			region.NameRu = raw["name_ru"].(string)
			region.Id = uint(pyraconv.ToInt64(raw["id"]))
		}
	}

	return
}

/**
 Get Location info by ip
 */
func (v *Visitor) GetLocation(ip string) (region model.Location, err error) {

	mmRecord, err := v.GetMaxMind(ip); if err == nil {
		region.Longitude = pyraconv.ToFloat32(mmRecord.Location.Longitude)
		region.Latitude = pyraconv.ToFloat32(mmRecord.Location.Latitude)
		region.TimeZone = mmRecord.Location.TimeZone
	}

	if err !=nil || mmRecord.Country.IsoCode == "RU" || mmRecord.Country.IsoCode == "UA" {

		spRecord, err := v.GetSyPexGeo(ip); if err == nil {
			raw := spRecord["city"].(map[string]interface{})
			region.Longitude = pyraconv.ToFloat32(raw["lon"])
			region.Latitude = pyraconv.ToFloat32(raw["lat"])
		}
	}

	return
}

/**
 Get Postal info by ip
 */
func (v *Visitor) GetPostal(ip string) (postal model.Postal, err error) {

	mmRecord, err := v.GetMaxMind(ip); if err == nil {
		postal.Code = mmRecord.Postal.Code
	}

	return
}

/**
 Get browser info by user agent and BrowsCap
 */
func (v *Visitor) GetBrowser(ua string) (browser model.Browser, err error) {

	total := time.Now()

	record, err := v.GetBrowsCap(ua); if err == nil {
		browser = model.Browser{
			Name:    record.Browser,
			Version: record.BrowserVersion,
			Type:    record.BrowserType,
		}
	}

	if v.Debug {
		fmt.Println("GetBrowser: Execute time " + time.Since(total).String())
	}
	return
}

/**
 Get platform info by user agent and BrowsCap
 */
func (v *Visitor) GetPlatform(ua string) (platform model.Platform, err error) {

	total := time.Now()

	record, err := v.GetBrowsCap(ua); if err == nil {
		platform = model.Platform{
			Name:    record.Platform,
			Version: record.PlatformVersion,
			Short:   record.PlatformShort,
		}
	}

	if v.Debug {
		fmt.Println("GetPlatform: Execute time " + time.Since(total).String())
	}
	return
}

/**
 Get device information by BrowsCap and UaSurfer
 */
func (v *Visitor) GetDevice(ua string) (device model.Device, err error) {
	total := time.Now()

	record, err := v.GetBrowsCap(ua); if err == nil {
		device = model.Device{
			Name:    record.DeviceName,
			Brand: 	 record.DeviceBrand,
			Type:    record.DeviceType,
		}
	}

	if device.Name == "unknown" || device.Name == "" {
		UaSurf := uasurfer.Parse(ua)

		switch info := UaSurf.DeviceType.String(); info {
		case "DeviceTV":
			device.Name = "Smart TV"
			device.Type = "TV Device"
			break
		case "DevicePhone":
			device.Name = "unknown"
			device.Type = "Mobile Phone"
			break
		case "DeviceComputer":
			device.Name = "unknown"
			device.Type = "Desktop"
			break
		case "DeviceTablet":
			device.Name = "unknown"
			device.Type = "TABLET"
			break
		case "DeviceConsole":
			device.Name = "unknown"
			device.Type = "TV Device"
			break
		case "DeviceWearable":
			device.Name = "unknown"
			device.Type = "Console"
			break
		}
	}

	if device.Type == "" {
		err = errors.New("Device not found " + ua)
		return
	}

	if v.Debug {
		fmt.Println("GetDevice: Execute time " + time.Since(total).String())
	}
	return
}

/**
 Get personal data
 */
func (v *Visitor) GetPersonal(ua string) (personal model.Personal, err error) {

	total := time.Now()

	personal = model.Personal{
		Ua: ua,
	}

	if v.Debug {
		fmt.Println("GetPlatform: Execute time " + time.Since(total).String())
	}

	return
}

/**
 Get user ip data
 */
func (v *Visitor) GetIp(ips string) (ip model.Ip, err error) {

	total := time.Now()

	ip = model.Ip{
		V4:ips,
	}

	if v.Debug {
		fmt.Println("GetIp: Execute time " + time.Since(total).String())
	}

	return
}