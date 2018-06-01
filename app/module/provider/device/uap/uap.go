package uap

import (
	"github.com/krecu/go-visitor/app/module/provider/device"
	"github.com/ua-parser/uap-go/uaparser"
)

type Uap struct {
	conn *uaparser.Parser
}

type Option struct {
	db string
}

func New(opt Option) (proto *Uap, err error) {

	proto = &Uap{}
	proto.conn, err = uaparser.New(opt.db)
	return
}

func (uap *Uap) Get(ua string) (proto *device.Model, err error) {
	data := uap.conn.Parse(ua)

	if err == nil {
		proto = &device.Model{
			Device: struct {
				Name  string
				Type  string
				Brand string
			}{
				Name:  data.Device.Model,
				Type:  "",
				Brand: data.Device.Brand,
			},
			Browser: struct {
				Name    string
				Type    string
				Version string
			}{
				Name:    data.UserAgent.Patch,
				Type:    "",
				Version: data.UserAgent.Minor,
			},
			Platform: struct {
				Name    string
				Short   string
				Version string
			}{
				Name:    data.Os.Patch,
				Short:   "",
				Version: data.Os.Minor,
			},
		}
	}
	return
}

func (uap *Uap) Close() {
	uap.Close()
}