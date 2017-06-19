package cache

import (
	"github.com/aerospike/aerospike-client-go"
	"errors"
	"github.com/CossackPyra/pyraconv"
	"github.com/krecu/go-visitor/model"
)

type AeroSpike struct {
	Host string
	Port int
	Ns string
	Db string
	Timeout int
	Client *aerospike.Client
}

func New(host string, port int, ns string, db string, timeout int) (*AeroSpike) {

	conn, err := aerospike.NewClient(host, port); if err != nil {
		panic(err)
	}

	return &AeroSpike{
		Host: host,
		Port: port,
		Db: db,
		Ns: ns,
		Timeout: timeout,
		Client: conn,
	}
}

/**
 Get value from cache
 */
func (c *AeroSpike) Get(id string) (vModel model.Visitor, err error) {

	var record *aerospike.Record

	policy := new(aerospike.BasePolicy)
	policy.Priority = aerospike.HIGH

	key, err := aerospike.NewKey(c.Ns, c.Db, id); if err != nil {
		return
	}

	record, err = c.Client.Get(policy, key); if record == nil {
		err = errors.New("Empty value")
	}

	if err == nil {
		vModel = c.UnMarshal(record.Bins)
	}

	return
}

func (c *AeroSpike) Set(visitor model.Visitor, extra map[string]interface{}) (err error) {

	record := c.Marshal(visitor)

	policy := new(aerospike.WritePolicy)
	policy.Priority = aerospike.HIGH
	policy.Expiration = 2592000 // month

	for key, val := range extra {
		record[key] = val
	}

	key, err := aerospike.NewKey(c.Ns, c.Db, record["id"]); if err != nil {
		return
	}

	err = c.Client.Put(policy, key, record); if err != nil {
		return
	}

	return
}


//
func (c *AeroSpike) Marshal (visitor model.Visitor) (map[string]interface{}) {

	record := make(map[string]interface{})

	// browser
	record["br_min"] 	= visitor.Browser.MinorVer
	record["br_maj"] 	= visitor.Browser.MajorVer
	record["br_type"] 	= visitor.Browser.Type
	record["br_ver"] 	= visitor.Browser.Version
	record["br_name"] 	= visitor.Browser.Name

	// device
	record["dc_name"] 	= visitor.Device.Name
	record["dc_type"] 	= visitor.Device.Type
	record["dc_brand"] 	= visitor.Device.Brand

	// platform
	record["pf_name"] 	= visitor.Platform.Name
	record["pf_short"] 	= visitor.Platform.Short
	record["pf_ver"] 	= visitor.Platform.Version
	record["pf_desc"] 	= visitor.Platform.Description
	record["pf_maker"] 	= visitor.Platform.Maker

	// city
	record["ct_id"] 	= visitor.City.Id
	record["ct_name"] 	= visitor.City.Name
	record["ct_name_ru"]	 = visitor.City.NameRu

	// country
	record["cn_id"] 	= visitor.Country.Id
	record["cn_name"] 	= visitor.Country.Name
	record["cn_name_ru"] 	= visitor.Country.NameRu
	record["cn_iso"] 	= visitor.Country.Iso

	// location
	record["lc_lat"] 	= visitor.Location.Latitude
	record["lc_lon"] 	= visitor.Location.Longitude
	record["lc_tz"] 	= visitor.Location.TimeZone

	// personal
	record["pr_ua"] 	= visitor.Personal.Ua
	record["pr_fn"] 	= visitor.Personal.FirstName
	record["pr_ln"] 	= visitor.Personal.LastName
	record["pr_pa"] 	= visitor.Personal.Patronymic
	record["pr_age"] 	= visitor.Personal.Age
	record["pr_ge"] 	= visitor.Personal.Gender

	// region
	record["re_id"] 	= visitor.Region.Id
	record["re_name"] 	= visitor.Region.Name
	record["re_name_ru"] 	= visitor.Region.NameRu

	// postal
	record["po_code"] 	= visitor.Postal.Code

	// ip
	record["ip_v4"] 	= visitor.Ip.V4
	record["ip_v6"] 	= visitor.Ip.V6

	return record
}

//
func (c *AeroSpike) UnMarshal(values map[string]interface{}) (visitor model.Visitor){

	var ok bool

	visitor.Browser = model.Browser{
		MinorVer: 	pyraconv.ToString(values["br_min"]),
		MajorVer: 	pyraconv.ToString(values["br_maj"]),
		Type: 		pyraconv.ToString(values["br_type"]),
		Version: 	pyraconv.ToString(values["br_ver"]),
		Name: 		pyraconv.ToString(values["br_name"]),
	}


	visitor.Device = model.Device{
		Brand: 		pyraconv.ToString(values["dc_brand"]),
		Type: 		pyraconv.ToString(values["dc_type"]),
		Name: 		pyraconv.ToString(values["dc_name"]),
	}

	visitor.Platform = model.Platform{
		Name: 		pyraconv.ToString(values["pf_name"]),
		Short: 		pyraconv.ToString(values["pf_short"]),
		Version: 	pyraconv.ToString(values["pf_ver"]),
		Description: 	pyraconv.ToString(values["pf_desc"]),
		Maker: 		pyraconv.ToString(values["pf_maker"]),
	}

	// if region not empty
	_, ok = values["re_id"]; if ok {
		visitor.Region = model.Region{
			Id: 	uint(pyraconv.ToInt64(values["re_id"])),
			Name:   pyraconv.ToString(values["re_name"]),
			NameRu: pyraconv.ToString(values["re_name_ru"]),
		}
	}

	// if country not empty
	_, ok = values["ct_id"]; if ok {
		visitor.City = model.City{
			Id: 	uint(pyraconv.ToInt64(values["ct_id"])),
			Name: 	pyraconv.ToString(values["ct_name"]),
			NameRu: pyraconv.ToString(values["ct_name_ru"]),
		}
	}

	// if city not empty
	_, ok = values["cn_id"]; if ok {
		visitor.City = model.City{
			Id: 	uint(pyraconv.ToInt64(values["cn_id"])),
			Name: 	pyraconv.ToString(values["cn_name"]),
			NameRu: pyraconv.ToString(values["cn_name_ru"]),
		}
	}

	visitor.Postal = model.Postal{
		Code:   	pyraconv.ToString(values["po_code"]),
	}

	visitor.Location = model.Location{
		Latitude: 	pyraconv.ToFloat32(values["lc_lat"]),
		Longitude:  	pyraconv.ToFloat32(values["lc_lon"]),
		TimeZone:   	pyraconv.ToString(values["lc_tz"]),
	}

	visitor.Ip = model.Ip{
		V4: 		pyraconv.ToString(values["ip_v4"]),
		V6:  		pyraconv.ToString(values["ip_v6"]),
	}

	visitor.Personal = model.Personal{
		Gender: 	pyraconv.ToString(values["pr_ge"]),
		Age:  		pyraconv.ToString(values["pr_age"]),
		Patronymic:   	pyraconv.ToString(values["pr_pa"]),
		LastName:   	pyraconv.ToString(values["pr_ln"]),
		FirstName:   	pyraconv.ToString(values["pr_fn"]),
		Ua:   		pyraconv.ToString(values["pr_ua"]),
	}

	return
}