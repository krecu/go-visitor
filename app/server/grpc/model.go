package grpc

import (
	"encoding/json"

	"github.com/CossackPyra/pyraconv"
	"github.com/krecu/go-visitor/app/module/wrapper"
)

const FieldCreated = "created"
const FieldUpdated = "updated"
const FieldId = "id"
const FieldExtra = "extra"
const FieldBrowserMinorVer = "br_min"
const FieldBrowserMajorVer = "br_maj"
const FieldBrowserType = "br_type"
const FieldBrowserVersion = "br_ver"
const FieldBrowserName = "br_name"
const FieldDeviceName = "dc_name"
const FieldDeviceType = "dc_type"
const FieldDeviceBrand = "dc_brand"
const FieldPlatformName = "pf_name"
const FieldPlatformShort = "pf_short"
const FieldPlatformVersion = "pf_ver"
const FieldPlatformDescription = "pf_desc"
const FieldPlatformMaker = "pf_maker"
const FieldCityId = "ct_id"
const FieldCityName = "ct_name"
const FieldCityNameRu = "ct_name_ru"
const FieldCountryId = "cn_id"
const FieldCountryName = "cn_name"
const FieldCountryNameRu = "cn_name_ru"
const FieldCountryIso = "cn_iso"
const FieldLocationLatitude = "lc_lat"
const FieldLocationLongitude = "lc_lon"
const FieldLocationTimeZone = "lc_tz"
const FieldPersonalUa = "pr_ua"
const FieldPersonalFirstName = "pr_fn"
const FieldPersonalLastName = "pr_ln"
const FieldPersonalPatronymic = "pr_pa"
const FieldPersonalAge = "pr_age"
const FieldPersonalGender = "pr_ge"
const FieldRegionId = "re_id"
const FieldRegionName = "re_name"
const FieldRegionNameRu = "re_name_ru"
const FieldPostalCode = "po_code"
const FieldIpV4 = "ip_v4"
const FieldIpV6 = "ip_v6"

type Record struct {
	Id    string                 `json:"id"`
	Ip    string                 `json:"ip"`
	Ua    string                 `json:"ua"`
	Extra map[string]interface{} `json:"extra"`
}

type Model struct {
	Id      string `json:"id"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
	City    struct {
		wrapper.City
		Mapping int `json:"mapping_id"`
	} `json:"city"`
	Country struct {
		wrapper.Country
		Mapping int `json:"mapping_id"`
	} `json:"country"`
	Region struct {
		wrapper.Region
		Mapping int `json:"mapping_id"`
	} `json:"region"`
	Browser struct {
		wrapper.Browser
		Mapping int `json:"mapping_id"`
	} `json:"browser"`
	Device struct {
		wrapper.Device
		Mapping int `json:"mapping_id"`
	} `json:"device"`
	Platform struct {
		wrapper.Platform
		Mapping int `json:"mapping_id"`
	} `json:"platform"`
	Personal struct {
		wrapper.Personal
		Mapping int `json:"mapping_id"`
	} `json:"personal"`
	Ip struct {
		wrapper.Ip
	} `json:"ip"`
	Postal struct {
		wrapper.Postal
	} `json:"postal"`
	Location struct {
		wrapper.Location
	} `json:"location"`
	Debug struct {
		wrapper.Debug
	} `json:"debug"`
	Extra map[string]interface{} `json:"extra"`
}

func (m *Model) Formed(info *wrapper.Model, data map[string]interface{}) (err error) {
	m.City.Id = info.Country.Id
	m.City.Name = info.Country.Name
	m.City.Mapping = 1

	m.Country.Id = info.Country.Id
	m.Country.Name = info.Country.Name
	m.Country.Iso = info.Country.Iso
	m.Country.Iso3166_1_alpha_3 = info.Country.Iso3166_1_alpha_3
	m.Country.Mapping = 1

	m.Region.Id = info.Region.Id
	m.Region.Name = info.Region.Name
	m.Region.Iso = info.Region.Iso
	m.Region.Mapping = 1

	m.Browser.Name = info.Browser.Name
	m.Browser.Type = info.Browser.Type
	m.Browser.Version = info.Browser.Version
	m.Browser.Mapping = 1

	m.Device.Name = info.Device.Name
	m.Device.Type = info.Device.Type
	m.Device.Brand = info.Device.Brand
	m.Device.Mapping = 1

	m.Platform.Name = info.Platform.Name
	m.Platform.Short = info.Platform.Short
	m.Platform.Version = info.Platform.Version
	m.Platform.Mapping = 1

	m.Personal.Ua = info.Personal.Ua

	m.Ip.V4 = info.Ip.V4

	m.Postal.Code = info.Postal.Code

	m.Location.Latitude = info.Location.Latitude
	m.Location.Longitude = info.Location.Longitude
	m.Location.TimeZone = info.Location.TimeZone

	m.Debug.TimeDevice = info.Debug.TimeDevice
	m.Debug.TimeGeo = info.Debug.TimeGeo
	m.Debug.ProviderDevice = info.Debug.ProviderDevice
	m.Debug.ProviderGeo = info.Debug.ProviderGeo

	// декомпозируем доп данные
	if data != nil {
		err = json.Unmarshal([]byte(pyraconv.ToString(data[FieldExtra])), &m.Extra)
	}

	return
}

func (m *Model) UnMarshal() (data map[string]interface{}, err error) {

	data = map[string]interface{}{
		FieldId:         m.Id,
		FieldCreated:    m.Created,
		FieldUpdated:    m.Updated,
		FieldPersonalUa: m.Personal.Ua,
		FieldIpV4:       m.Ip.V4,
		FieldExtra:      m.Extra,
	}

	data[FieldExtra], err = json.Marshal(m.Extra)

	return
}
func (m *Model) Marshal() (buf []byte, err error) {

	buf, err = json.Marshal(m)

	return
}
