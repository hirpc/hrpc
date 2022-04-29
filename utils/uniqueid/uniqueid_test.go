package uniqueid

import (
	"fmt"
	"testing"
)

func TestUniqueid(t *testing.T) {
	t.Log(String())
}

func BenchmarkUniqueid(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fmt.Println(String())
	}
}
