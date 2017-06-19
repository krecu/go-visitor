package server

import (
	"github.com/krecu/go-visitor/cache"
	"github.com/krecu/go-visitor/model"
	"github.com/krecu/go-visitor"
)

type Core struct {
	Debug bool
	Db string
	Buffer int
	Cache *cache.AeroSpike
	Core *visitor.Visitor
}

func New(core *visitor.Visitor, cache *cache.AeroSpike) (*Core) {
	return &Core{
		Core: core,
		Cache: cache,
	}
}

func (c *Core) Get(id string, ip string, ua string, extra map[string]interface{}) (info model.Visitor, err error) {


	info, err = c.Cache.Get(id); if err != nil {

		info, err = c.Core.Identify(ip, ua)

		go func() {

			c.Cache.Set(info, extra)

		}()
	}

	return
}

