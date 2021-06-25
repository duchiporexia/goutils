package xmsg

import "testing"

func BenchmarkUIntId(b *testing.B) {
	id := LongID(599696599696)
	for i := 0; i < b.N; i++ {
		id.String()
	}
}
