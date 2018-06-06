package wrapper

import (
	"time"
)

type Debug struct {
	TimeGeo        time.Duration
	ProviderGeo    string
	TimeDevice     time.Duration
	ProviderDevice string
	TimeTotal      time.Duration
}

type City struct {
	Name string `json:"name"`
	Id   string `json:"geoname_id"`
}

type Country struct {
	Name              string `json:"name"`
	Id                string `json:"geoname_id"`
	Iso               string `json:"iso_code"`
	Iso3166_1_alpha_3 string `json:"iso_code_3166_1_alpha_3"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TimeZone  string  `json:"time_zone"`
}

type Postal struct {
	Code string `json:"code"`
}

type Region struct {
	Name string `json:"name"`
	Id   string `json:"geoname_id"`
	Iso  string `json:"iso"`
}

type Geo struct {
	City     City     `json:"city"`
	Country  Country  `json:"country"`
	Location Location `json:"location"`
	Region   Region   `json:"region"`
	Postal   Postal   `json:"postal"`
}

type Browser struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type Device struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Brand string `json:"brand"`
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
	Name    string `json:"name"`
	Short   string `json:"short"`
	Version string `json:"version"`
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
