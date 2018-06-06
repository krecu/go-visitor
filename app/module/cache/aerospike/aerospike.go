package aerospike

import (
	"strings"

	"time"

	"fmt"

	"reflect"

	"github.com/CossackPyra/pyraconv"
	"github.com/aerospike/aerospike-client-go"
)

type Option struct {
	Hosts             []string
	ConnectionTimeout time.Duration
	GetTimeout        time.Duration
	WriteTimeout      time.Duration
	Expiration        uint32
	Set               string
	NameSpace         string
}

type Cache struct {
	opt  Option
	conn *aerospike.Client
}

func New(opt Option) (proto *Cache, err error) {

	proto = &Cache{
		opt: opt,
	}

	hosts := []*aerospike.Host{}

	for _, h := range opt.Hosts {
		arg := strings.Split(h, ":")
		hosts = append(hosts, &aerospike.Host{
			Name: pyraconv.ToString(arg[0]),
			Port: int(pyraconv.ToInt64(arg[1])),
		})
	}

	policy := aerospike.NewClientPolicy()
	policy.Timeout = opt.ConnectionTimeout * time.Millisecond

	proto.conn, err = aerospike.NewClientWithPolicyAndHost(policy, hosts...)

	return
}

// добавляем в кеш
func (c *Cache) Set(key string, value interface{}) (err error) {

	var (
		cacheKey *aerospike.Key
	)

	policy := aerospike.NewWritePolicy(0, c.opt.Expiration)
	policy.Priority = aerospike.HIGH
	policy.Timeout = c.opt.WriteTimeout * time.Millisecond

	cacheKey, err = aerospike.NewKey(
		c.opt.NameSpace,
		c.opt.Set,
		key,
	)
	if err != nil {
		return
	}

	if record, ok := value.(map[string]interface{}); ok {
		err = c.conn.Put(policy, cacheKey, record)
	} else {
		err = fmt.Errorf("Неверный формат")
	}

	return
}

// добавляем в кеш
func (c *Cache) SetExpired(key string, value interface{}) (err error) {
	return
}

// извлекаем из кеша
func (c *Cache) Get(key string, value interface{}) (err error) {

	var (
		cacheKey *aerospike.Key
	)

	policy := aerospike.NewPolicy()
	policy.Priority = aerospike.HIGH
	policy.Timeout = c.opt.GetTimeout * time.Millisecond

	cacheKey, err = aerospike.NewKey(
		c.opt.NameSpace,
		c.opt.Set,
		key,
	)
	if err != nil {
		return
	}

	rv := reflect.ValueOf(value).Elem()
	info, err := c.conn.Get(policy, cacheKey)
	if err == nil {
		for k, v := range info.Bins {
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
	}

	return
}

func (c *Cache) List(prefix string) (items []interface{}, err error) {

	return
}

func (c *Cache) Del(key string) {
}

func (c *Cache) Close() {
}
