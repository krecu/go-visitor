package model

type Geo struct {
	City City		`json:"city"`
	Country Country		`json:"country"`
	Location Location	`json:"location"`
	Region Region		`json:"region"`
	Postal Postal		`json:"postal"`
}
