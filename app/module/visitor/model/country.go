package model

type Country struct {
	Name              string `json:"name"`
	NameRu            string `json:"name_ru"`
	Id                uint   `json:"geoname_id"`
	Iso               string `json:"iso_code"`
	Iso3166_1_alpha_3 string `json:"iso_code_3166_1_alpha_3"`
	Mapping           int    `json:"mapping_id"`
}
