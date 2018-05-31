package model

import (
	"fmt"

	"github.com/CossackPyra/pyraconv"
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
	Updated  int64                  `json:"updated"`
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

	// дата создания
	record[FieldCreated] = visitor.Created
	// дата обновления
	record[FieldUpdated] = visitor.Updated
	// уникальный идентификатор юзера
	record[FieldId] = visitor.Id
	// добавляем доп поля
	record[FieldExtra] = visitor.Extra
	// personal
	record[FieldPersonalUa] = visitor.Personal.Ua
	// ip
	record[FieldIpV4] = visitor.Ip.V4

	return record
}

func VisitorUnMarshal(values map[string]interface{}) (proto *Visitor) {

	proto = &Visitor{}

	proto.Created = pyraconv.ToInt64(values[FieldCreated])
	proto.Updated = pyraconv.ToInt64(values[FieldUpdated])
	proto.Id = pyraconv.ToString(values[FieldId])
	proto.Ip = Ip{
		V4: pyraconv.ToString(values[FieldIpV4]),
	}

	proto.Personal = Personal{
		Ua: pyraconv.ToString(values[FieldPersonalUa]),
	}

	if extra, ok := values[FieldExtra].(map[interface{}]interface{}); ok {
		proto.Extra = make(map[string]interface{})
		for key, value := range extra {
			if k, ok := key.(string); ok {
				proto.Extra[k] = cleanupMapValue(value)
			}
		}
	}

	return
}

func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = cleanupMapValue(v)
	}
	return res
}

func cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
