package model

import (
	"github.com/CossackPyra/pyraconv"
)

const FieldCreated = "created"
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

type Debug struct {
	GeoProvider    string
	DeviceProvider string
	TimeGeo        string
	TimeDevice     string
	TimeTotal      string
}

type Visitor struct {
	Id       string                 `json:"id"`
	Created  int64                  `json:"created"`
	City     City                   `json:"city"`
	Country  Country                `json:"country"`
	Location Location               `json:"location"`
	Postal   Postal                 `json:"postal"`
	Region   Region                 `json:"region"`
	Browser  Browser                `json:"browser"`
	Device   Device                 `json:"device"`
	Platform Platform               `json:"platform"`
	Personal Personal               `json:"personal"`
	Ip       Ip                     `json:"ip"`
	Extra    map[string]interface{} `json:"extra"`
	Debug    Debug                  `json:"debug"`
}

func VisitorMarshal(visitor *Visitor) map[string]interface{} {

	record := make(map[string]interface{})

	record[FieldCreated] = visitor.Created
	record[FieldId] = visitor.Id

	// добавляем доп поля
	record[FieldExtra] = visitor.Extra

	// browser
	record[FieldBrowserMinorVer] = visitor.Browser.MinorVer
	record[FieldBrowserMajorVer] = visitor.Browser.MajorVer
	record[FieldBrowserType] = visitor.Browser.Type
	record[FieldBrowserVersion] = visitor.Browser.Version
	record[FieldBrowserName] = visitor.Browser.Name

	// device
	record[FieldDeviceName] = visitor.Device.Name
	record[FieldDeviceType] = visitor.Device.Type
	record[FieldDeviceBrand] = visitor.Device.Brand

	// platform
	record[FieldPlatformName] = visitor.Platform.Name
	record[FieldPlatformShort] = visitor.Platform.Short
	record[FieldPlatformVersion] = visitor.Platform.Version
	record[FieldPlatformDescription] = visitor.Platform.Description
	record[FieldPlatformMaker] = visitor.Platform.Maker

	// city
	record[FieldCityId] = visitor.City.Id
	record[FieldCityName] = visitor.City.Name
	record[FieldCityNameRu] = visitor.City.NameRu

	// country
	record[FieldCountryId] = visitor.Country.Id
	record[FieldCountryName] = visitor.Country.Name
	record[FieldCountryNameRu] = visitor.Country.NameRu
	record[FieldCountryIso] = visitor.Country.Iso

	// location
	record[FieldLocationLatitude] = visitor.Location.Latitude
	record[FieldLocationLongitude] = visitor.Location.Longitude
	record[FieldLocationTimeZone] = visitor.Location.TimeZone

	// personal
	record[FieldPersonalUa] = visitor.Personal.Ua
	record[FieldPersonalFirstName] = visitor.Personal.FirstName
	record[FieldPersonalLastName] = visitor.Personal.LastName
	record[FieldPersonalPatronymic] = visitor.Personal.Patronymic
	record[FieldPersonalAge] = visitor.Personal.Age
	record[FieldPersonalGender] = visitor.Personal.Gender

	// region
	record[FieldRegionId] = visitor.Region.Id
	record[FieldRegionName] = visitor.Region.Name
	record[FieldRegionNameRu] = visitor.Region.NameRu

	// postal
	record[FieldPostalCode] = visitor.Postal.Code

	// ip
	record[FieldIpV4] = visitor.Ip.V4
	record[FieldIpV6] = visitor.Ip.V6

	return record
}

func VisitorUnMarshal(values map[string]interface{}) (proto *Visitor) {

	proto = &Visitor{}

	proto.Created = pyraconv.ToInt64(values[FieldCreated])
	proto.Id = pyraconv.ToString(values[FieldId])

	proto.Browser = Browser{
		MinorVer: pyraconv.ToString(values[FieldBrowserMinorVer]),
		MajorVer: pyraconv.ToString(values[FieldBrowserMajorVer]),
		Type:     pyraconv.ToString(values[FieldBrowserType]),
		Version:  pyraconv.ToString(values[FieldBrowserVersion]),
		Name:     pyraconv.ToString(values[FieldBrowserName]),
	}

	proto.Device = Device{
		Brand: pyraconv.ToString(values[FieldDeviceBrand]),
		Type:  pyraconv.ToString(values[FieldDeviceType]),
		Name:  pyraconv.ToString(values[FieldDeviceName]),
	}

	proto.Platform = Platform{
		Name:        pyraconv.ToString(values[FieldPlatformName]),
		Short:       pyraconv.ToString(values[FieldPlatformShort]),
		Version:     pyraconv.ToString(values[FieldPlatformVersion]),
		Description: pyraconv.ToString(values[FieldPlatformDescription]),
		Maker:       pyraconv.ToString(values[FieldPlatformMaker]),
	}

	proto.Device.Mapping = DeviceMapping(proto.Device, proto.Platform)

	if region, ok := values[FieldRegionId]; ok {
		proto.Region = Region{
			Id:      uint(pyraconv.ToInt64(region)),
			Name:    pyraconv.ToString(values[FieldRegionName]),
			NameRu:  pyraconv.ToString(values[FieldRegionNameRu]),
			Mapping: RegionMapping(pyraconv.ToString(values[FieldRegionNameRu])),
		}
	}

	if city, ok := values[FieldCityId]; ok {
		proto.City = City{
			Id:      uint(pyraconv.ToInt64(city)),
			Name:    pyraconv.ToString(values[FieldCityName]),
			NameRu:  pyraconv.ToString(values[FieldCityNameRu]),
			Mapping: CityMapping(pyraconv.ToString(values[FieldCityNameRu])),
		}
	}

	if country, ok := values[FieldCountryId]; ok {
		proto.Country = Country{
			Id:                uint(pyraconv.ToInt64(country)),
			Name:              pyraconv.ToString(values[FieldCountryName]),
			NameRu:            pyraconv.ToString(values[FieldCountryNameRu]),
			Iso:               pyraconv.ToString(values[FieldCountryIso]),
			Iso3166_1_alpha_3: ISO_3166_1_alpha_3Mapping(pyraconv.ToString(values[FieldCountryIso])),
			Mapping:           CountryMapping(pyraconv.ToString(values[FieldCountryNameRu])),
		}
	}

	proto.Postal = Postal{
		Code: pyraconv.ToString(values[FieldPostalCode]),
	}

	proto.Location = Location{
		Latitude:  pyraconv.ToFloat32(values[FieldLocationLatitude]),
		Longitude: pyraconv.ToFloat32(values[FieldLocationLongitude]),
		TimeZone:  pyraconv.ToString(values[FieldLocationTimeZone]),
	}

	proto.Ip = Ip{
		V4: pyraconv.ToString(values[FieldIpV4]),
		V6: pyraconv.ToString(values[FieldIpV6]),
	}

	proto.Personal = Personal{
		Gender:     pyraconv.ToString(values[FieldPersonalGender]),
		Age:        pyraconv.ToString(values[FieldPersonalAge]),
		Patronymic: pyraconv.ToString(values[FieldPersonalPatronymic]),
		LastName:   pyraconv.ToString(values[FieldPersonalLastName]),
		FirstName:  pyraconv.ToString(values[FieldPersonalFirstName]),
		Ua:         pyraconv.ToString(values[FieldPersonalUa]),
	}

	if extra, ok := values[FieldExtra].(map[interface{}]interface{}); ok {
		proto.Extra = make(map[string]interface{})
		for key, value := range extra {
			if k, ok := key.(string); ok {
				proto.Extra[k] = value
			}
		}
	}

	return
}
