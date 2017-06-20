package cache

import "github.com/krecu/go-visitor/model"

type Cache interface {
	Get(id string) (visitor model.Visitor, err error)
	Set(id string, visitor model.Visitor) (err error)
}