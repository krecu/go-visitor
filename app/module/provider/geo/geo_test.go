package geo

import (
	"testing"

	"github.com/k0kubun/pp"
)

func TestLoadCountry(t *testing.T) {

	list, err := LoadCountry()
	if err != nil {
		t.Error(err)
	} else {
		pp.Println(list)
	}
}
