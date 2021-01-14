package neighborhood

import "testing"

// benchmark outputs to avoid compiler optimizations
var idx Index
var result []Point

func BenchmarkLoad_1k(b *testing.B) {
	points := globalPoints(1_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx = NewIndex().Load(points...)
	}
}

func BenchmarkLoad_10k(b *testing.B) {
	points := globalPoints(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx = NewIndex().Load(points...)
	}
}

func BenchmarkLoad_100k(b *testing.B) {
	points := globalPoints(100_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx = NewIndex().Load(points...)
	}
}

func BenchmarkNearby_100k_k1(b *testing.B) {
	benchmarkNearby(b, 100_000,1)
}

func BenchmarkNearby_100k_k10(b *testing.B) {
	benchmarkNearby(b, 100_000,10)
}

func BenchmarkNearby_100k_k100(b *testing.B) {
	benchmarkNearby(b, 100_000,100)
}


func benchmarkNearby(b *testing.B, n, k int) {
	points := globalPoints(n)
	origin := namedPoint("seattle")
	idx := NewIndex().Load(points...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = idx.Nearby(origin, k, AcceptAny)
	}
}