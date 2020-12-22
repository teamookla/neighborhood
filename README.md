# neighborhood
A very fast, static, in-memory nearest neighbor (KNN) search index for locations on Earth.
Accounts for Earth's curvature and date line wrapping. Utilizes a [k-d tree](https://en.wikipedia.org/wiki/K-d_tree)
for very quick spatial searches.

* :white_check_mark: Zero dependency
* :white_check_mark: Pure Go
* :white_check_mark: 97.6% test coverage
* :white_check_mark: Used in production by high-volume services powering [Speedtest](https://speedtest.net)

## Usage
### Install
```
go get github.com/teamookla/neighborhood
```

### Import
```go
import github.com/teamookla/neighborhood
```

### Implement the Point interface
Index your own custom types by implementing the simple `Point` interface.
```go
func (t *Thing) Lat() float64 { return t.latitude }
func (t *Thing) Lon() float64 { return t.longitude }
```
Alternately, you can add `Coordinates` to your type, which implements `Point` interface.
```go
type Thing struct {
	neighborhood.Coordinates,
	...
}
```

### Create an index
Create a new index with all searchable `Points`. Indexes are static and immutable.
```go
idx := neighborhood.NewIndex(things...)
```

### Search for `k` Nearest Neighbors
```go
origin := neighborhood.NewCoordinates(-122, 47) // origin can be any Point
results := idx.Nearby(origin, k, neighborhood.AcceptAny)
```

### Custom `Accepter` function (optional)
You can specify criteria other than distance that `Points` must meet to be included in results. 
```go
accepter := func (p Point) bool {
	return p.(*Thing).Color == "Blue" // we only want Blue things
}
origin := neighborhood.NewCoordinates(-122, 47) // origin can be any Point
results := idx.Nearby(origin, k, accepter)
```

### Implement `Ranker` (optional)
You can optionally specify a secondary search rank for a `Point` (distance is the primary).
```go
// prefer older Things if there are multiple at the same distance
func (t *Thing) GetRank() float64 { return float64(t.age) }
```

## Performance
Benchmark tests get k-nearest-neighbors from an index of 260,281 Points (spread evenly around the globe).
Tests were run on a 2019 Macbook Pro 16.
```
BenchmarkNearby_k1-16              57852             20140 ns/op
BenchmarkNearby_k10-16             34261             34629 ns/op
BenchmarkNearby_k100-16            15967             75622 ns/op
```

## Attribution
Neighborhood was inspired by Mapbox Engineer [Vladimir Agafonkin's](https://github.com/mourner) excellent 
[dive into spatial search algorithms](https://blog.mapbox.com/a-dive-into-spatial-search-algorithms-ebd0c5e39d2a), 
and  it is based on his [geokdbush](https://github.com/mourner/geokdbush) javascript library.
