package neighborhood

import "testing"

func BenchmarkNearby_k1(b *testing.B) {
	benchmarkNearby(b, 1)
}

func BenchmarkNearby_k10(b *testing.B) {
	benchmarkNearby(b, 10)
}

func BenchmarkNearby_k100(b *testing.B) {
	benchmarkNearby(b, 100)
}


func benchmarkNearby(b *testing.B, k int) {
	points := globalPoints()
	origin := namedPoint("seattle")
	idx := NewIndex(points...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Nearby(origin, k, AcceptAny)
	}
}