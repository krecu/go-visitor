package uasurfer

import (
	"fmt"

	"github.com/avct/uasurfer"
	"github.com/krecu/go-visitor/app/module/provider/device"
)

type UaSurfer struct {
	opt Option
}

type Option struct {
	Weight int
	Name   string
}

func New(opt Option) (proto *UaSurfer, err error) {
	proto = &UaSurfer{
		opt: opt,
	}
	return
}

func (uas *UaSurfer) Weight() int {
	return uas.opt.Weight
}

func (uas *UaSurfer) Name() string {
	return uas.opt.Name
}

func (uas *UaSurfer) Get(ua string) (proto *device.Model, err error) {
	data := uasurfer.Parse(ua)
	proto = &device.Model{}

	if data != nil {
		switch data.DeviceType.String() {
		case "DeviceUnknown":
			proto.Device.Name = "unknown"
			proto.Device.Brand = "unknown"
			proto.Device.Type = "unknown"
			break
		case "DeviceTV":
			proto.Device.Name = "Smart TV"
			proto.Device.Brand = "unknown"
			proto.Device.Type = "TV Device"
			break
		case "DevicePhone":
			proto.Device.Name = "unknown"
			proto.Device.Brand = "unknown"
			proto.Device.Type = "Mobile Phone"
			break
		case "DeviceComputer":
			proto.Device.Name = "unknown"
			proto.Device.Brand = "unknown"
			proto.Device.Type = "Desktop"
			break
		case "DeviceTablet":
			proto.Device.Name = "unknown"
			proto.Device.Brand = "unknown"
			proto.Device.Type = "Tablet"
			break
		case "DeviceConsole":
			proto.Device.Name = "unknown"
			proto.Device.Brand = "unknown"
			proto.Device.Type = "Console"
			break
		case "DeviceWearable":
			proto.Device.Name = "unknown"
			proto.Device.Brand = "unknown"
			proto.Device.Type = "Mobile Device"
			break
		}

		switch data.OS.Name.String() {
		case "OSUnknown":
			proto.Platform.Name = "unknown"
		case "OSWindowsPhone":
			proto.Platform.Name = "WinPhone"
		case "OSWindows":
			proto.Platform.Name = "Win32"
		case "OSMacOSX":
			proto.Platform.Name = "macOS"
			proto.Device.Brand = "Apple"
			proto.Device.Name = "Macintosh"
		case "OSiOS":
			proto.Platform.Name = "iOS"
			proto.Device.Brand = "Apple"
		case "OSAndroid":
			proto.Platform.Name = "Android"
		case "OSBlackberry":
			proto.Platform.Name = "Blackberry"
		case "OSChromeOS":
			proto.Platform.Name = "Chrome"
		case "OSKindle":
			proto.Platform.Name = "Kindle"
		case "OSWebOS":
			proto.Platform.Name = "WebOS"
		case "OSLinux":
			proto.Platform.Name = "Linux"
		case "OSPlaystation":
			proto.Platform.Name = "OrbisOS"
			proto.Device.Name = "Playstation"
			proto.Device.Brand = "Sony"
		case "OSXbox":
			proto.Platform.Name = "Xbox OS"
			proto.Device.Name = "Xbox One"
			proto.Device.Brand = "Microsoft"
		case "OSNintendo":
			proto.Platform.Name = "Nintendo"
			proto.Device.Name = "Switch"
			proto.Device.Brand = "Nintendo"
		case "OSBot":
			proto.Platform.Name = "Bot"
		}

		switch data.OS.Platform.String() {
		case "PlatformUnknown":
			proto.Platform.Short = "unknown"
		case "PlatformWindows":
			proto.Platform.Short = "Windows"
		case "PlatformMac":
			proto.Platform.Short = "mac"
		case "PlatformLinux":
			proto.Platform.Short = "Linux"
		case "PlatformiPad":
			proto.Platform.Short = "iPad"
		case "PlatformiPhone":
			proto.Platform.Short = "iPhone"
		case "PlatformiPod":
			proto.Platform.Short = "iPod"
		case "PlatformBlackberry":
			proto.Platform.Short = "Blackberry"
		case "PlatformWindowsPhone":
			proto.Platform.Short = "WindowsPhone"
		case "PlatformPlaystation":
			proto.Platform.Short = "Playstation"
		case "PlatformXbox":
			proto.Platform.Short = "Xbox"
		case "PlatformNintendo":
			proto.Platform.Short = "Nintendo"
		case "PlatformBot":
			proto.Platform.Short = "Bot"
		}

		switch data.Browser.Name.String() {
		case "BrowserUnknown":
			proto.Browser.Name = "unknown"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserChrome":
			proto.Browser.Name = "Chrome"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserIE":
			proto.Browser.Name = "IE"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserSafari":
			proto.Browser.Name = "Safari"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserFirefox":
			proto.Browser.Name = "Firefox"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserAndroid":
			proto.Browser.Name = "Android"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserOpera":
			proto.Browser.Name = "Opera"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserBlackberry":
			proto.Browser.Name = "Blackberry"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserUCBrowser":
			proto.Browser.Name = "UCBrowser"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserSilk":
			proto.Browser.Name = "Silk"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserNokia":
			proto.Browser.Name = "Nokia"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserNetFront":
			proto.Browser.Name = "NetFront"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserQQ":
			proto.Browser.Name = "QQ"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserMaxthon":
			proto.Browser.Name = "Maxthon"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserSogouExplorer":
			proto.Browser.Name = "SogouExplorer"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserSpotify":
			proto.Browser.Name = "Spotify"
			proto.Browser.Type = "unknown"
			proto.Browser.Version = "unknown"
		case "BrowserBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserAppleBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserBaiduBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserBingBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserDuckDuckGoBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserFacebookBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserGoogleBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserLinkedInBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserMsnBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserPingdomBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserTwitterBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserYandexBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		case "BrowserYahooBot":
			proto.Browser.Name = "Bot"
			proto.Browser.Type = "Bot"
			proto.Browser.Version = "Bot"
		}

	}

	if proto.Device.Type == "unknown" {
		err = fmt.Errorf("Bad device type")
	}

	return
}

func (uas *UaSurfer) Close() {
}
