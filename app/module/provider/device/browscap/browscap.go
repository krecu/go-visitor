package browscap

import (
	"time"

	"github.com/krecu/browscap_go"
	"github.com/krecu/go-visitor/app/module/provider/device"
)

type BrowsCap struct {
	opt Option
}

type Option struct {
	db string
}

func New(opt Option) (proto *BrowsCap, err error) {

	proto = &BrowsCap{
		opt: opt,
	}
	err = browscap_go.InitBrowsCap(
		opt.db,
		true,
		time.Duration(3600)*time.Second,
		time.Duration(3600)*time.Second,
	)

	return
}

func (bc *BrowsCap) Get(ua string) (proto *device.Model, err error) {
	var (
		browser *browscap_go.Browser
	)
	browser, err = browscap_go.GetBrowser(ua)

	if err == nil {
		proto = &device.Model{
			Device: struct {
				Name  string
				Type  string
				Brand string
			}{
				Name:  browser.DeviceName,
				Type:  browser.DeviceType,
				Brand: browser.DeviceBrand,
			},
			Browser: struct {
				Name    string
				Type    string
				Version string
			}{
				Name:    browser.Browser,
				Type:    browser.BrowserType,
				Version: browser.BrowserVersion,
			},
			Platform: struct {
				Name    string
				Short   string
				Version string
			}{
				Name:    browser.Platform,
				Short:   browser.PlatformShort,
				Version: browser.PlatformVersion,
			},
		}
	}

	return
}

func (bc *BrowsCap) Close() {
}
