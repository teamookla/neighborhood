# neighborhood
A very fast, static, in-memory nearest neighbor (KNN) search index for locations on Earth.
Approximates Earth's curvature and accounts for International Date Line wrapping using
[Haversine great-circle distance formula](https://en.wikipedia.org/wiki/Haversine_formula). Utilizes a [k-d tree](https://en.wikipedia.org/wiki/K-d_tree)
for very quick spatial searches.

* :white_check_mark: Zero dependency
* :white_check_mark: Pure Go
* :white_check_mark: 100% test coverage
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

### Create an Index and load Points
Create a new index and load all searchable `Points`.
Each call to `Load` will replace all `Points` in the `Index` with the provided `Points`.
```go
idx := neighborhood.NewIndex().Load(things...)
```
`Load` mutates and returns the `Index` to allow call chaining.
If you don't use call chaining, the returned `Index` can be ignored.
```go
idx := neighborhood.NewIndex()
idx.Load(things...)
```

Each call to `Add` will take the provided points, merge them with the current points in the index.
`Add` uses `Load` behind the scenes, so it is only as performant as calling `Load`.
`Add` supports call chaining like `Load` does as well.
```go
idx := neighborhood.NewIndex()
idx.Load(things...)
idx.Add(aFewMoreThings...)
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
Benchmark tests get k-nearest-neighbors from an Index with default options and 100,000 Points
(uniformly distributed around the globe). Tests were run on a 2019 Macbook Pro 16.

```
BenchmarkLoad_1k-16                        50382             24706 ns/op
BenchmarkLoad_10k-16                        1821            649872 ns/op
BenchmarkLoad_100k-16                        168           6915047 ns/op
BenchmarkNearby_100k_k1-16                 72946             16338 ns/op
BenchmarkNearby_100k_k10-16                71701             16658 ns/op
BenchmarkNearby_100k_k100-16               12033             97936 ns/op
```

## Attribution
Neighborhood was inspired by Mapbox Engineer [Vladimir Agafonkin's](https://github.com/mourner) excellent 
[dive into spatial search algorithms](https://blog.mapbox.com/a-dive-into-spatial-search-algorithms-ebd0c5e39d2a), 
and  it is based on his [geokdbush](https://github.com/mourner/geokdbush) javascript library.
