package model

type Region struct {
	Name    string `json:"name"`
	NameRu  string `json:"name_ru"`
	Id      uint   `json:"geoname_id"`
	Iso     string `json:"iso"`
	Mapping int    `json:"mapping_id"`
}
