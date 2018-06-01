package wrapper

import (
	"time"
)

type Debug struct {
	TimeGeo    time.Duration
	TimeDevice time.Duration
	TimeTotal  time.Duration
}

type City struct {
	Name    string `json:"name"`
	NameRu  string `json:"name_ru"`
	Id      uint   `json:"geoname_id"`
	Mapping int    `json:"mapping_id"`
}

type Country struct {
	Name              string `json:"name"`
	NameRu            string `json:"name_ru"`
	Id                uint   `json:"geoname_id"`
	Iso               string `json:"iso_code"`
	Iso3166_1_alpha_3 string `json:"iso_code_3166_1_alpha_3"`
	Mapping           int    `json:"mapping_id"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	TimeZone  string  `json:"time_zone"`
}

type Postal struct {
	Code string `json:"code"`
}

type Region struct {
	Name    string `json:"name"`
	NameRu  string `json:"name_ru"`
	Id      uint   `json:"geoname_id"`
	Iso     string `json:"iso"`
	Mapping int    `json:"mapping_id"`
}

type Geo struct {
	City     City     `json:"city"`
	Country  Country  `json:"country"`
	Location Location `json:"location"`
	Region   Region   `json:"region"`
	Postal   Postal   `json:"postal"`
}

type Browser struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Version  string `json:"version"`
	MajorVer string `json:"majorver"`
	MinorVer string `json:"minorver"`
	Mapping  int    `json:"mapping_id"`
}

type Device struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Brand   string `json:"brand"`
	Mapping int    `json:"mapping_id"`
}

type Ip struct {
	V4 string `json:"v4"`
	V6 string `json:"v6"`
}

type Personal struct {
	Ua         string `json:"ua"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	Age        string `json:"age"`
	Gender     string `json:"gender"`
}

type Platform struct {
	Name        string `json:"name"`
	Short       string `json:"short"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Maker       string `json:"maker"`
	Mapping     int    `json:"mapping_id"`
}

type System struct {
	Device   Device   `json:"device"`
	Browser  Browser  `json:"browser"`
	Platform Platform `json:"platform"`
}

type Model struct {
	City     City     `json:"city"`
	Country  Country  `json:"country"`
	Location Location `json:"location"`
	Postal   Postal   `json:"postal"`
	Region   Region   `json:"region"`
	Browser  Browser  `json:"browser"`
	Device   Device   `json:"device"`
	Platform Platform `json:"platform"`
	Personal Personal `json:"personal"`
	Ip       Ip       `json:"ip"`
	Debug    Debug    `json:"debug"`
}
