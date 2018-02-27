package main

import (
	"time"

	"fmt"

	"net"

	"encoding/json"

	"github.com/CossackPyra/pyraconv"
	"github.com/aerospike/aerospike-client-go"
	"github.com/avct/uasurfer"
	"github.com/krecu/browscap_go"
	"github.com/krecu/go-visitor/model"
	"github.com/oschwald/geoip2-golang"
)

var VisitorErrorEmpty = fmt.Errorf("Visitor empty ")

type VisitorService struct {
	app *App
}

func NewVisitorService(app *App) (proto *VisitorService, err error) {

	proto = &VisitorService{
		app: app,
	}

	return
}

// получение
func (v *VisitorService) Get(id string) (proto *model.Visitor, err error) {

	var (
		record *aerospike.Record
		key    *aerospike.Key
	)

	// общий счетчик времени
	_total := time.Now()

	if !v.app.aerospike.IsConnected() {
		err = fmt.Errorf("Нет соединения с хранилищем ")
		return
	}

	policy := aerospike.NewPolicy()
	policy.Priority = aerospike.HIGH
	policy.Timeout = v.app.config.GetDuration("app.aerospike.GetTimeout") * time.Millisecond

	key, err = aerospike.NewKey(
		v.app.config.GetString("app.aerospike.NameSpace"),
		v.app.config.GetString("app.aerospike.Set"),
		id,
	)
	if err != nil {
		return
	}

	record, err = v.app.aerospike.Get(policy, key)
	if record == nil || err != nil {
		if err == nil {
			err = VisitorErrorEmpty
		}
		return
	}

	proto = model.VisitorUnMarshal(record.Bins)

	err = v.Indent(proto)

	if err != nil {
		return
	}

	// считаем время для дебага
	proto.Debug.TimeTotal = time.Since(_total).String()

	return
}

// создание
func (v *VisitorService) Post(id string, ip string, ua string, extra map[string]interface{}) (proto *model.Visitor, err error) {

	// первоначальная модель
	proto = &model.Visitor{
		Id:      id,
		Created: time.Now().Unix(),
		Updated: time.Now().Unix(),
		Extra:   extra,
		Ip: model.Ip{
			V4: ip,
		},
		Personal: model.Personal{
			Ua: ua,
		},
	}

	err = v.Indent(proto)

	if err != nil {
		return
	}

	// сораняем в БД
	v.Save(
		v.app.aerospike,
		proto,
		v.app.config.GetString("app.aerospike.NameSpace"),
		v.app.config.GetString("app.aerospike.Set"),
		v.app.config.GetDuration("app.aerospike.WriteTimeout")*time.Millisecond,
		uint32(v.app.config.GetInt("app.aerospike.Ttl")),
	)

	return
}

// изменение
func (v *VisitorService) Patch(id string, fields map[string]interface{}) (proto *model.Visitor, err error) {

	proto, err = v.Get(id)

	if err != nil {
		return
	} else {
		if proto == nil {
			err = VisitorErrorEmpty
			return
		}
	}

	proto.Updated = time.Now().Unix()

	if _, ok := fields["extra"]; ok {
		var buf []byte
		buf, err = json.Marshal(fields["extra"])
		if err == nil {
			err = json.Unmarshal(buf, &proto.Extra)
		}
	}

	if err == nil {
		v.Save(
			v.app.aerospike,
			proto,
			v.app.config.GetString("app.aerospike.NameSpace"),
			v.app.config.GetString("app.aerospike.Set"),
			v.app.config.GetDuration("app.aerospike.WriteTimeout")*time.Millisecond,
			uint32(v.app.config.GetInt("app.aerospike.Ttl")),
		)
	}

	return
}

// удаление
func (v *VisitorService) Delete(id string) (err error) {

	var (
		key *aerospike.Key
	)

	if !v.app.aerospike.IsConnected() {
		err = fmt.Errorf("Нет соединения с хранилищем ")
		return
	}

	policy := aerospike.NewWritePolicy(1, 1)
	policy.Priority = aerospike.HIGH
	policy.Timeout = v.app.config.GetDuration("app.aerospike.GetTimeout") * time.Millisecond

	key, err = aerospike.NewKey(
		v.app.config.GetString("app.aerospike.NameSpace"),
		v.app.config.GetString("app.aerospike.Set"),
		id,
	)
	if err != nil {
		return
	}

	_, err = v.app.aerospike.Delete(policy, key)
	return
}

// определение информации о юзере
// применяется в GET/POST/PATCH
func (v *VisitorService) Indent(proto *model.Visitor) (err error) {

	var (
		DataGeo    *model.Geo
		DataDevice *model.System
		TimerTotal = time.Now()
	)

	// счетчик времени для определения гео
	TimerGeo := time.Now()

	// считываем ГЕО данные по SyPex базе
	if DataGeo, err = v.GetGeo(proto.Ip.V4); err == nil {
		proto.Country = DataGeo.Country
		proto.Region = DataGeo.Region
		proto.City = DataGeo.City
		proto.Postal = DataGeo.Postal
		proto.Location = DataGeo.Location
	} else {
		Logger.Errorf("Ошибка определения географии: IP:%s; ERR:%s", proto.Ip.V4, err)
	}

	// считаем время для дебага
	proto.Debug.TimeGeo = time.Since(TimerGeo).String()

	// счетчик времени для определения устройств
	TimerDevice := time.Now()

	// считываем ГЕО данные по SyPex базе
	if DataDevice, err = v.GetSystem(proto.Personal.Ua); err == nil {
		proto.Device = DataDevice.Device
		proto.Browser = DataDevice.Browser
		proto.Platform = DataDevice.Platform
	} else {
		Logger.Errorf("Ошибка определения устройства: UA:%s; ERR:%s", proto.Personal.Ua, err)
	}

	// считаем время для дебага
	proto.Debug.TimeDevice = time.Since(TimerDevice).String()

	// считаем время для дебага
	proto.Debug.TimeTotal = time.Since(TimerTotal).String()

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

// определение информации об устройстве
func (v *VisitorService) GetSystem(ua string) (device *model.System, err error) {

	var (
		cacheKey = fmt.Sprintf("device_%s", ua)
	)

	device = &model.System{}

	if GeoCache, ok := v.app.cache.Get(cacheKey); GeoCache != nil && ok {
		if value, ok := GeoCache.(*model.System); ok {
			device = value
			return
		}
	}

	// считываем данные по устройству/браузеру/платформе
	if DeviceData, ok := browscap_go.GetBrowser(ua); ok {

		device.Device = model.Device{
			Name:  DeviceData.DeviceName,
			Brand: DeviceData.DeviceBrand,
			Type:  DeviceData.DeviceType,
		}

		device.Platform = model.Platform{
			Name:    DeviceData.Platform,
			Version: DeviceData.PlatformVersion,
			Short:   DeviceData.PlatformShort,
		}

		device.Browser = model.Browser{
			Name:    DeviceData.Browser,
			Version: DeviceData.BrowserVersion,
			Type:    DeviceData.BrowserType,
		}

		// если в основной базе не найденно значений
		if DeviceData.DeviceType == "unknown" || DeviceData.DeviceType == "" {

			// Дополнительный парсер, если вдруг основной сбойнул
			// работает по мапингам:
			// Mobile Phone, Mobile Device, Tablet, Desktop, TV Device, Console,
			// FonePad, Ebook Reader, Car Entertainment System or unknown
			if DeviceTmp := uasurfer.Parse(ua); DeviceData != nil {
				switch DeviceTmp.DeviceType.String() {
				case "DeviceTV":
					device.Device.Name = "Smart TV"
					device.Device.Type = "TV Device"
					break
				case "DevicePhone":
					device.Device.Name = "unknown"
					device.Device.Type = "Mobile Phone"
					break
				case "DeviceComputer":
					device.Device.Name = "unknown"
					device.Device.Type = "Desktop"
					break
				case "DeviceTablet":
					device.Device.Name = "unknown"
					device.Device.Type = "Tablet"
					break
				case "DeviceConsole":
					device.Device.Name = "unknown"
					device.Device.Type = "Console"
					break
				case "DeviceWearable":
					device.Device.Name = "unknown"
					device.Device.Type = "Mobile Device"
					break
				}

			}
		}

		device.Device.Mapping = model.DeviceMapping(device.Device, device.Platform)
	} else {
		err = fmt.Errorf("Not found user agent: %s ", ua)
	}

	if err == nil && device != nil {
		v.app.cache.SetDefault(cacheKey, device)

	}

	return
}

// сохранение
func (v *VisitorService) Save(cache *aerospike.Client, proto *model.Visitor, ns string, set string, timeout time.Duration, ttl uint32) (err error) {

	var (
		key *aerospike.Key
	)

	if !cache.IsConnected() {
		err = fmt.Errorf("Нет соединения с хранилищем ")
		return
	}

	// преобразуем структуру в массив
	record := model.VisitorMarshal(proto)

	policy := aerospike.NewWritePolicy(0, aerospike.TTLServerDefault)
	policy.Priority = aerospike.HIGH
	policy.Timeout = timeout
	if ttl == 0 {
		policy.Expiration = aerospike.TTLDontExpire
	} else {
		policy.Expiration = ttl
	}

	key, err = aerospike.NewKey(
		ns,
		set,
		proto.Id,
	)
	if err != nil {
		Logger.Error(err)
		return
	}

	err = cache.Put(policy, key, record)
	if err != nil {
		Logger.Error(err)
	}

	return
}
