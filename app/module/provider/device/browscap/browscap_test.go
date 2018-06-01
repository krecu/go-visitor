package browscap

import (
	"testing"

	"github.com/k0kubun/pp"
)

func TestBrowsCap_Get(t *testing.T) {
	bc, err := New(Option{
		db: "./full_php_browscap.ini",
	})

	if err != nil {
		t.Error(err)
	}

	list := []string{
		"Android/4.4.4 stb/Eltex NV501WAC.NV501WAC/NV501WAC.NV501WAC.armeabi-v7a",
		"Peers.TV/1.0.1 Linux/3.3.8-3.1 stb/Eltex NV100/mips",
		"Android/7.1.2 stb/Amlogic Amlogic.X96mini/p281.p281.armeabi-v7a",
		"Android/4.2.2 stb/Eltex unknown.NV310WAC/bcm7429.NV310WAC.mips",
		"Android/4.4.4 stb/Eltex NV501WAC.NV501WAC/NV501WAC.armeabi-v7a",
		"Android/6.0.1 stb/Amlogic Amlogic.MBOX/p212.p212.armeabi-v7a",
		"Android/4.2.2 stb/Eltex unknown.NV310WAC/bcm7429.mips",
		"Android/4.4.4 stb/Eltex NV501.NV501/NV501.NV501.armeabi-v7a",
		"Android/5.1.1 stb/rockchip Android.rk322x-box/rk30sdk.rk322x_box.armeabi-v7a",
		"Android/7.1.2 stb/Amlogic Amlogic.TX3_Mini/p281.p281.armeabi-v7a",
		"AndroidTV/6.0.1 stb/Xiaomi Xiaomi.MIBOX3/once.once.armeabi-v7a",
		"Android/6.0.1 stb/Amlogic Amlogic.X92/q201.q201.armeabi-v7",
		"Android/4.4.4 stb/CVTE changhong.CH-HW338-DTV-00-00/CH-HW338-DTV-00-00_PB803.CH-HW338-DTV-00-00_PB803.armeabi-v7a",
		"Android/5.1.1 stb/amlogic Android.T95N-2G/p201_tn2.p201_tn2.armeabi-v7a",
	}

	for _, s := range list {
		d, err := bc.Get(s)
		pp.Println(d, err)
	}
}

func BenchmarkBrowsCap_Get(b *testing.B) {
	bc, err := New(Option{
		db: "./full_php_browscap.ini",
	})

	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		bc.Get("Mozilla/5.0 (iPhone; CPU iPhone OS 11_2_6 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) Mobile/15D100 Instagram 28.0.0.12.285 (iPhone9,3; iOS 11_2_6; ru_RU; ru-RU; scale=2.00; gamut=wide; 750x1334)")
	}
}
