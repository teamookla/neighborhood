package neighborhood

import "testing"

// benchmark outputs to avoid compiler optimizations
var idx Index
var result []Point

func BenchmarkNewIndex_1k(b *testing.B) {
	points := globalPoints()[0:1000]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx = NewIndex(points...)
	}
}

func BenchmarkNewIndex_10k(b *testing.B) {
	points := globalPoints()[0:10000]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx = NewIndex(points...)
	}
}

func BenchmarkNewIndex_100k(b *testing.B) {
	points := globalPoints()[0:100000]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx = NewIndex(points...)
	}
}

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
		result = idx.Nearby(origin, k, AcceptAny)
	}
}