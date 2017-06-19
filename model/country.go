package model

type Country struct {
	Name string	`json:"name"`
	NameRu string	`json:"name_ru"`
	Id uint		`json:"geoname_id"`
	Iso string	`json:"iso_code"`
}
