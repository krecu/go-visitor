package model

type Visitor struct {
	City City		`json:"city"`
	Country Country		`json:"country"`
	Location Location	`json:"location"`
	Postal Postal		`json:"postal"`
	Region Region		`json:"region"`
	Browser Browser		`json:"browser"`
	Device Device		`json:"device"`
	Platform Platform	`json:"platform"`
	Personal Personal	`json:"personal"`
	Ip Ip			`json:"ip"`
}