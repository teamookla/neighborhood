// Neighborhood is a very fast, static, in-memory nearest neighbor (KNN) search index for locations on Earth.
// Accounts for Earth's curvature and date line wrapping. Utilizes a k-d tree for very quick spatial searches.
package neighborhood

// Index interface defines the nearest-neighbor search contract
type Index interface {
	// Nearby finds the k nearest Points to the origin that meet the Accepter criteria.
	// If there are multiple Points that are the same distance from the origin and the Points implement the Ranker
	// interface, the higher ranking Points will be preferred. Nearby may return less than k results if it cannot
	// find k Points in the Index that meet the Accepter criteria.
	Nearby(p Point, k int, accept Accepter) []Point

	// Load adds searchable Points to the Index.
	// Each call to Load will replace all Points in the Index with the provided Points.
	// Load mutates and returns the Index to allow call chaining.
	Load(points ...Point) Index
}

// Point interface defines latitude and longitude accessors
type Point interface {
	// Lat gets Point latitude
	Lat() float64
	// Lon gets Point longitude
	Lon() float64
}

// Ranker is an optional interface to define a secondary Point rank (distance is the primary)
type Ranker interface {
	// GetRank gets Point rank (secondary sorting property)
	GetRank() float64
}

// Accepter defines a function that will accept or ignore a given Point
type Accepter func(p Point) bool

// AcceptAny is a simple Accepter implementation that accepts any and all Points
func AcceptAny(Point) bool { return true }

// NewIndex creates a new neighborhood Index with default options
func NewIndex() Index {
	// uses a kd-tree index by default
	return NewKDTreeIndex(DefaultKDTreeOptions())
}

