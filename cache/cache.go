package cache

type Cache interface {
	Get(id string) (data map[string]interface{}, err error)
	Set(data map[string]interface{}) bool
}