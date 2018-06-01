package device

type Device interface {
	Get(ua string) (proto *Model, err error)
	Close()
}
