/**
	Провайдер хранения данных о пользователе в AeroSpike реалезованные по интерфейсу Cache
 */

package cache

import (
	"github.com/aerospike/aerospike-client-go"
	"errors"
	"github.com/CossackPyra/pyraconv"
	"github.com/krecu/go-visitor/model"
	_ "log"
	"fmt"
	"log"
	"time"
)

type AeroSpike struct {
	Host string
	Port int
	Ns string
	Db string
	Timeout time.Duration
	Ttl uint32
	Client *aerospike.Client
}

/**
	host - адрес сервера
	port - порт
	ns - name space базы
	db - set хранения данных
	timeout - глобальная задержка на все операции
 */
func New(host string, port int, ns string, db string, timeout time.Duration, ttl uint32) (*AeroSpike) {

	conn, err := aerospike.NewClient(host, port); if err != nil {
		panic(err)
	}

	return &AeroSpike{
		Host: host,
		Port: port,
		Db: db,
		Ns: ns,
		Timeout: timeout,
		Ttl: ttl, // @todo вынести в конфиг
		Client: conn,
	}
}

/**
	Закрываем соединение
 */
func (c *AeroSpike) Close() {
	c.Client.Close()
}

/**
 	id - идентификатор в бд
 */
func (c *AeroSpike) Get(id string) (visitor model.Visitor, err error) {

	if !c.Client.IsConnected() {
		log.Fatalf("AeroSpike disconect")
	}

	var record *aerospike.Record

	policy := new(aerospike.BasePolicy)
	policy.Priority = aerospike.HIGH
	policy.Timeout = c.Timeout * time.Millisecond

	key, err := aerospike.NewKey(c.Ns, c.Db, id); if err != nil {
		return
	}

	record, err = c.Client.Get(policy, key); if record == nil {
		err = errors.New("Empty value")
	}

	if err == nil {
		visitor = c.UnMarshal(record.Bins)
	}

	return
}

/**
	visitor - модель от ядра
	extra - дополнительный набор полей
 */
func (c *AeroSpike) Set(id string, visitor model.Visitor) (err error) {

	if !c.Client.IsConnected() {
		log.Fatalf("AeroSpike disconect")
	}

	// преобразуем структуру в массив
	record := c.Marshal(visitor)

	policy := new(aerospike.WritePolicy)
	policy.Priority = aerospike.HIGH
	policy.Expiration = c.Ttl
	policy.Timeout = c.Timeout * time.Millisecond

	// генерируем digits
	key, err := aerospike.NewKey(c.Ns, c.Db, id); if err != nil {
		return
	}

	err = c.Client.Put(policy, key, record); if err != nil {
		return
	}

	return
}

/**
	Упаковываем модель в key/value массив. Так как в AeroSpike есть несколько ограничения
	на длину имя bin (14 символов) то нам необходимо сократить имена
 */
func (c *AeroSpike) Marshal (visitor model.Visitor) (map[string]interface{}) {

	record := make(map[string]interface{})

	record["created"] 	= visitor.Created
	record["id"] 		= visitor.Id

	// добавляем доп поля
	record["extra"] 	= visitor.Extra

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

/**
	Распаковываем массив вида key/value значений в модель ядра
 */
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
		visitor.Country = model.Country{
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

	visitor.Created = pyraconv.ToInt64(values["created"])
	visitor.Id = pyraconv.ToString(values["id"])

	_, ok = values["extra"]; if ok {
		visitor.Extra = _cleanupInterfaceMap(values["extra"].(map[interface{}]interface{}))
	}

	return
}

/**
	Helper function - преобразуем интерфейс в key/value
 */
func _cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = _cleanupMapValue(v)
	}
	return res
}

/**
	Helper function - преобразуем массив интерфейсов в key/value
 */
func _cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = _cleanupMapValue(v)
	}
	return res
}

/**
	Helper function - преобразуем массив интерфейс/интерфейсов в key/value
 */
func _cleanupMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return _cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return _cleanupInterfaceMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
