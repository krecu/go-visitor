package device

type Device interface {
	Get(ua string) (proto *Model, err error)
	Weight() int
	Name() string
	Close()
}

type OrderProvider []Device

func (a OrderProvider) Len() int      { return len(a) }
func (a OrderProvider) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a OrderProvider) Less(i, j int) bool {
	return a[i].Weight() > a[j].Weight()
}
