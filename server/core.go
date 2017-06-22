/**

	Основной проксик между ядром visitora и клиентами/серверами
	Состоит из:
	1) Get - метод получения данных о посетителе
		- Пытается получить данные из кеша
		- Если неудалось запрашивает новые данные
		- Возвращает key/value значения
 */

package server

import (
	"github.com/krecu/go-visitor/cache"
	"github.com/krecu/go-visitor/model"
	"github.com/krecu/go-visitor"
	_ "log"
	"time"
)

type Core struct {
	Debug bool
	Db string
	Buffer int
	Cache cache.Cache
	Core *visitor.Visitor
}

/**
	core - движок опредеоения данных о пользователе
	cache - интерфейс бд для кеширования/ временного хранения информации о пользователе
 */
func New(core *visitor.Visitor, cache cache.Cache) (*Core) {
	return &Core{
		Core: core,
		Cache: cache,
	}
}

/**
	Определение визитора и запись в бд

	id - идентификатор в бд
	ip - v4 IP адрес пользователя
	ua - user agent пользователя
	extra - дополнительные поля пользователя каторые не вычисляются а просто запишутся в бд
 */
func (c *Core) Get(id string, ip string, ua string, extra map[string]interface{}) (info model.Visitor, err error) {

	info, err = c.Cache.Get(id)

	// валидируем что в кеше записано то что нужно
	// иначе прсто перезапишем данные
	if info.City.Name == ""   || info.Country.Name == ""  || err != nil ||
	   	info.Device.Type == "" || info.Platform.Name == "" || info.Browser.Name == "" {

		info, err = c.Core.Identify(ip, ua)

		info.Id = id
		info.Created = time.Now().Unix()
		info.Extra = extra

		// обновляем данные в бд
		c.Cache.Set(id, info)
	}

	return
}

/**
	Обновление полей или запись в бд при остутствии

	id - идентификатор в бд
	ip - v4 IP адрес пользователя
	ua - user agent пользователя
	extra - дополнительные поля пользователя каторые не вычисляются а просто запишутся в бд
 */
func (c *Core) Put(id string, ip string, ua string, extra map[string]interface{}) (info model.Visitor, err error) {

	info, err = c.Cache.Get(id)

	if info.City.Name == ""   || info.Country.Name == ""  || err != nil ||
		info.Device.Type == "" || info.Platform.Name == "" || info.Browser.Name == "" {
		info, err = c.Core.Identify(ip, ua)
	}

	info.Id = id
	if info.Created == 0 {
		info.Created = time.Now().Unix()
	}

	info.Extra = extra

	// обновляем данные в бд
	c.Cache.Set(id, info)

	return
}

/**
	Полное обновение запись

	id - идентификатор в бд
	ip - v4 IP адрес пользователя
	ua - user agent пользователя
	extra - дополнительные поля пользователя каторые не вычисляются а просто запишутся в бд
 */
func (c *Core) Refresh(id string, ip string, ua string, extra map[string]interface{}) (info model.Visitor, err error) {


	info, err = c.Core.Identify(ip, ua)

	info.Id = id
	if info.Created == 0 {
		info.Created = time.Now().Unix()
	}

	info.Extra = extra

	// обновляем данные в бд
	c.Cache.Set(id, info)

	return
}

