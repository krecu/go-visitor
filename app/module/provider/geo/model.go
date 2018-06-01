package geo

type Model struct {
	Country struct {
		Id                 string
		Name               string
		Iso                string
		ISO_3166_1_alpha_3 string
	}
	Region struct {
		Id   string
		Name string
		Iso  string
	}
	City struct {
		Id   string
		Name string
	}
	Location struct {
		Latitude  float64
		Longitude float64
		TimeZone  string
	}
	Ip struct {
		V4 string
		V6 string
	}
	Postal struct {
		Code string
	}
}
