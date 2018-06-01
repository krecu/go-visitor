package videonow

import (
	"fmt"

	"github.com/krecu/go-visitor/app/module/cache"
	"github.com/krecu/go-visitor/app/module/cache/go-cache"
)

type Store struct {
	cache cache.Cache
	opt   Options
}

type Options struct {
	CacheExpire int64
}

func New(opt Options) (proto *Store, err error) {

	proto = &Store{
		opt: opt,
	}

	proto.cache, err = go_cache.New(86400)
	if err != nil {
		return
	}

	return
}

// получаем
func (s *Store) GetUser(id string) (proto string, err error) {

	var (
		cachePrefix = fmt.Sprintf("vn_user_%s", id)
	)

	err = s.cache.Get(cachePrefix, &proto)

	return
}

// запись
func (s *Store) SetUser(id string, proto string) (err error) {

	var (
		cachePrefix = fmt.Sprintf("vn_user_%s", id)
	)

	err = s.cache.SetExpired(cachePrefix, proto)

	return
}

// удаление
func (s *Store) DelUser(id string) {

	var (
		cachePrefix = fmt.Sprintf("vn_user_%s", id)
	)

	s.cache.Del(cachePrefix)

	return
}
