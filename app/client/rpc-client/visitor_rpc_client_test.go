package visitor_rpc_client

import (
	"testing"

	"github.com/k0kubun/pp"
)

func TestClient_Post(t *testing.T) {

	// соединение с визитором
	visitor, err := New([]string{"127.0.0.1:8094"})
	if err != nil {
		t.Errorf("%s", err)
	}

	_, err = visitor.Post("test", "79.104.42.249", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36", map[string]interface{}{
		"test": 1,
	})
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestClient_Get(t *testing.T) {

	// соединение с визитором
	visitor, _ := New([]string{"127.0.0.1:8094"})

	data, err := visitor.Get("test")
	if err != nil {
		t.Errorf("%s", err)
	}

	pp.Println(data)
}

func BenchmarkClient_Post(b *testing.B) {
	// соединение с визитором
	visitor, err := New([]string{"127.0.0.1:8094"})
	if err != nil {
		b.Errorf("%s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		visitor.Post("test", "79.104.42.249", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36", nil)
	}
}
