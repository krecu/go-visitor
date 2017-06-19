package model

type City struct {
	Name string	`json:"name"`
	NameRu string	`json:"name_ru"`
	Id uint		`json:"geoname_id"`
}
