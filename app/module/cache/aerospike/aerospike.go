package go_cache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-stats/app/module/cache"
	"time"

	vendor "github.com/patrickmn/go-cache"
)

type Cache struct {
	conn *vendor.Cache
}

func New(expireTtl int64) (proto *Cache, err error) {

	proto = &Cache{}
	proto.conn = vendor.New(time.Duration(expireTtl)*time.Second, time.Duration(expireTtl*2)*time.Second)

	return
}

// добавляем в кеш
func (c *Cache) Set(key string, value interface{}) (err error) {

	if buf, err := cache.Marshal(value); err == nil {
		c.conn.Set(key, buf, vendor.NoExpiration)
	} else {
		err = fmt.Errorf("Cache: %s", cache.NOT_FOUND)
	}

	return
}

// добавляем в кеш
func (c *Cache) SetExpired(key string, value interface{}) (err error) {

	if buf, err := cache.Marshal(value); err == nil {
		c.conn.SetDefault(key, buf)
	} else {
		err = fmt.Errorf("Cache: %s", cache.NOT_FOUND)
	}

	return
}

// извлекаем из кеша
func (c *Cache) Get(key string, value interface{}) (err error) {

	var (
		jsonData []byte
	)

	if buf, ok := c.conn.Get(key); ok {
		if jsonData, err = cache.Unmarshal(buf.([]byte)); err != nil {
			err = fmt.Errorf("Cache: %s, %s", cache.ERROR_UNPACK, err)
		} else {
			if err = json.Unmarshal(jsonData, &value); err != nil {
				err = fmt.Errorf("Cache: %s, %s", cache.ERROR_JSON, err)
			}
		}
	} else {
		err = fmt.Errorf("Cache: %s", cache.NOT_FOUND)
	}

	return
}

func (c *Cache) List(prefix string) (items []interface{}, err error) {

	var (
		buf       map[string]vendor.Item
		key       []byte
		prefixBuf = []byte(prefix)
	)

	buf = c.conn.Items()

	if len(buf) == 0 {
		err = cache.NOT_FOUND
	} else {
		for k, val := range buf {
			key = []byte(k)
			if !bytes.HasPrefix(key, prefixBuf) {
				continue
			}

			items = append(items, val.Object)
		}

	}

	return
}

func (c *Cache) Del(key string) {
	c.conn.Delete(key)
}

func (c *Cache) Close() {
	c.conn.Flush()
}
