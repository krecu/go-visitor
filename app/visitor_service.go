package main

import (
	"time"

	"fmt"

	"net"

	"encoding/json"

	"github.com/CossackPyra/pyraconv"
	"github.com/aerospike/aerospike-client-go"
	"github.com/avct/uasurfer"
	"github.com/imdario/mergo"
	"github.com/krecu/browscap_go"
	"github.com/krecu/go-visitor/model"
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

	// считаем время для дебага
	proto.Debug.TimeTotal = time.Since(_total).String()

	return
}

// создание
func (v *VisitorService) Post(id string, ip string, ua string, extra map[string]interface{}) (proto *model.Visitor, err error) {

	// общий счетчик времени
	_total := time.Now()

	proto = &model.Visitor{
		Id:      id,
		Created: time.Now().Unix(),
	}

	// счетчик времени для определения гео
	_geo := time.Now()

	// считываем ГЕО данные по SyPex базе
	if GeoData, ok := v.app.sypexgeo.GetCityFull(ip); ok == nil {

		proto.Debug.GeoProvider = "sypexgeo"

		if country, ok := GeoData["country"].(map[string]interface{}); ok {
			proto.Country.Id = uint(pyraconv.ToInt64(country["id"]))
			proto.Country.Name = pyraconv.ToString(country["name_en"])
			proto.Country.NameRu = pyraconv.ToString(country["name_ru"])
			proto.Country.Iso = pyraconv.ToString(country["iso"])
			proto.Country.Iso3166_1_alpha_3 = model.ISO_3166_1_alpha_3Mapping(pyraconv.ToString(country["iso"]))
			proto.Country.Mapping = model.CountryMapping(proto.Country.NameRu)
			proto.Location.TimeZone = pyraconv.ToString(country["timezone"])
		}
		if region, ok := GeoData["region"].(map[string]interface{}); ok {
			proto.Region.Id = uint(pyraconv.ToInt64(region["id"]))
			proto.Region.Name = pyraconv.ToString(region["name_en"])
			proto.Region.NameRu = pyraconv.ToString(region["name_ru"])
			proto.Region.Iso = pyraconv.ToString(region["iso"])
			proto.Region.Mapping = model.RegionMapping(proto.Region.NameRu)
		}
		if city, ok := GeoData["city"].(map[string]interface{}); ok {
			proto.City.Id = uint(pyraconv.ToInt64(city["id"]))
			proto.City.Name = pyraconv.ToString(city["name_en"])
			proto.City.NameRu = pyraconv.ToString(city["name_ru"])
			proto.Location.Latitude = pyraconv.ToFloat32(city["lat"])
			proto.Location.Longitude = pyraconv.ToFloat32(city["lon"])
			proto.City.Mapping = model.CityMapping(proto.City.NameRu)
		}
	}

	// дополянем и перероверяем ГЕО по MaxMind
	if GeoData, ok := v.app.maxmind.City(net.ParseIP(ip)); ok == nil {
		if country := GeoData.Country; country.GeoNameID != 0 && proto.Country.Id == 0 {
			proto.Debug.GeoProvider = "maxmind"
			proto.Country.Id = uint(pyraconv.ToInt64(country.GeoNameID))
			proto.Country.Name = pyraconv.ToString(country.Names["en"])
			proto.Country.NameRu = pyraconv.ToString(country.Names["ru"])
			proto.Country.Iso = pyraconv.ToString(country.IsoCode)
			proto.Country.Mapping = model.CountryMapping(country.Names["ru"])
		}
		if city := GeoData.City; city.GeoNameID != 0 && proto.City.Id == 0 {
			proto.City.Id = uint(pyraconv.ToInt64(city.GeoNameID))
			proto.City.Name = pyraconv.ToString(city.Names["en"])
			proto.City.NameRu = pyraconv.ToString(city.Names["ru"])
			proto.City.Mapping = model.CityMapping(city.Names["ru"])
		}
		if location := GeoData.Location; location.TimeZone != "" && proto.Location.TimeZone == "" {
			proto.Location.Latitude = pyraconv.ToFloat32(location.Latitude)
			proto.Location.Longitude = pyraconv.ToFloat32(location.Longitude)
			proto.Location.TimeZone = pyraconv.ToString(location.TimeZone)
		}
		if postal := GeoData.Postal; postal.Code != "" {
			proto.Postal.Code = pyraconv.ToString(postal.Code)
		}
	}

	// считаем время для дебага
	proto.Debug.TimeGeo = time.Since(_geo).String()

	// счетчик времени для определения устройств
	_device := time.Now()

	// считываем данные по устройству/браузеру/платформе
	if DeviceData, ok := browscap_go.GetBrowser(ua); ok {

		proto.Debug.DeviceProvider = "browscap"

		proto.Device = model.Device{
			Name:  DeviceData.DeviceName,
			Brand: DeviceData.DeviceBrand,
			Type:  DeviceData.DeviceType,
		}

		proto.Platform = model.Platform{
			Name:    DeviceData.Platform,
			Version: DeviceData.PlatformVersion,
			Short:   DeviceData.PlatformShort,
		}

		proto.Browser = model.Browser{
			Name:    DeviceData.Browser,
			Version: DeviceData.BrowserVersion,
			Type:    DeviceData.BrowserType,
		}

		// если в основной базе не найденно значений
		if DeviceData.DeviceType == "unknown" || DeviceData.DeviceType == "" {

			proto.Debug.DeviceProvider = "uasurfer"

			// Дополнительный парсер, если вдруг основной сбойнул
			// работает по мапингам:
			// Mobile Phone, Mobile Device, Tablet, Desktop, TV Device, Console,
			// FonePad, Ebook Reader, Car Entertainment System or unknown
			if DeviceTmp := uasurfer.Parse(ua); DeviceData != nil {
				switch DeviceTmp.DeviceType.String() {
				case "DeviceTV":
					proto.Device.Name = "Smart TV"
					proto.Device.Type = "TV Device"
					break
				case "DevicePhone":
					proto.Device.Name = "unknown"
					proto.Device.Type = "Mobile Phone"
					break
				case "DeviceComputer":
					proto.Device.Name = "unknown"
					proto.Device.Type = "Desktop"
					break
				case "DeviceTablet":
					proto.Device.Name = "unknown"
					proto.Device.Type = "Tablet"
					break
				case "DeviceConsole":
					proto.Device.Name = "unknown"
					proto.Device.Type = "Console"
					break
				case "DeviceWearable":
					proto.Device.Name = "unknown"
					proto.Device.Type = "Mobile Device"
					break
				}

			}
		}

		proto.Device.Mapping = model.DeviceMapping(proto.Device, proto.Platform)
	}

	// считаем время для дебага
	proto.Debug.TimeDevice = time.Since(_device).String()

	// по пользователю
	proto.Personal.Ua = ua
	proto.Ip.V4 = ip
	proto.Extra = extra

	// считаем время для дебага
	proto.Debug.TimeTotal = time.Since(_total).String()

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

	var mapProto map[string]interface{}

	proto, err = v.Get(id)
	if err != nil {
		return
	} else {
		if proto == nil {
			err = VisitorErrorEmpty
			return
		}
	}

	if bufProto, err := json.Marshal(proto); err == nil {
		if err = json.Unmarshal(bufProto, &mapProto); err == nil {
			if err = mergo.Merge(&fields, mapProto); err == nil {
				if bufProto, err := json.Marshal(fields); err == nil {
					if err = json.Unmarshal(bufProto, &proto); err == nil {
						go v.Save(
							v.app.aerospike,
							proto,
							v.app.config.GetString("app.aerospike.NameSpace"),
							v.app.config.GetString("app.aerospike.Set"),
							v.app.config.GetDuration("app.aerospike.WriteTimeout")*time.Millisecond,
							uint32(v.app.config.GetInt("app.aerospike.Ttl")),
						)
					}
				}
			}
		}
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

// удаление
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
