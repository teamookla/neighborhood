package neighborhood

import (
	"fmt"
	"testing"
)

var cities = namedPoints()
func Example() {
	// Create a new Index and load all searchable Points.
	idx := NewIndex().Load(cities...)
	// origin can be any Point
	origin := NewCoordinates(-122, 47)
	// find the 2 closest cities
	results := idx.Nearby(origin, 2, AcceptAny)
	fmt.Printf("%s is the closest city", results[0].(*NamedPoint).Name)
	fmt.Printf("%s is the second closest city", results[1].(*NamedPoint).Name)
}

func TestNewIndex_AnonymousCoordinates(t *testing.T) {
	type Thing struct{
		Coordinates
		name string
	}
	thing := &Thing{
		Coordinates: Coordinates{lon: -122.123, lat: 47.123},
		name:        "thing",
	}

	idx := NewIndex().Load(thing)
	origin := NewCoordinates(-122, 47)
	results := idx.Nearby(origin, 5, AcceptAny)
	assertEqual(t, 1, len(results))
	assertEqual(t, "thing", results[0].(*Thing).name)
}

func TestKDTree_Nearby_Simple(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex().Load(pts...)
	origin := NewCoordinates(-115, 45)

	results := idx.Nearby(origin, 3, AcceptAny)

	assertEqual(t, 3, len(results))
	assertEqual(t, "woodinville", results[0].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[1].(*NamedPoint).Name)
	assertEqual(t, "memphis", results[2].(*NamedPoint).Name)
}

func TestKDTree_Load_And_Add(t *testing.T) {
	pts := []Point{
		namedPoint("tokyo"),
		namedPoint("memphis"),
		namedPoint("cairo"),
		namedPoint("seattle"),
	}
	idx := NewIndex().Load(pts...)
	origin := NewCoordinates(-115, 45)

	results := idx.Nearby(origin, 4, AcceptAny)
	assertEqual(t, 4, len(results))
	assertEqual(t, "seattle", results[0].(*NamedPoint).Name)
	assertEqual(t, "memphis", results[1].(*NamedPoint).Name)
	assertEqual(t, "tokyo", results[2].(*NamedPoint).Name)
	assertEqual(t, "cairo", results[3].(*NamedPoint).Name)

	additionalPts := []Point{
		namedPoint("woodinville"),
		namedPoint("anchorage"),
		namedPoint("saopaulo"),
		namedPoint("eastrussia"),
	}

	// add the remaining points, validate that the new points work
	idx.Add(additionalPts...)
	results = idx.Nearby(origin, 6, AcceptAny)
	assertEqual(t, 6, len(results))
	assertEqual(t, "woodinville", results[0].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[1].(*NamedPoint).Name)
	assertEqual(t, "memphis", results[2].(*NamedPoint).Name)
	assertEqual(t, "anchorage", results[3].(*NamedPoint).Name)
	assertEqual(t, "eastrussia", results[4].(*NamedPoint).Name)
	assertEqual(t, "tokyo", results[5].(*NamedPoint).Name)

	// make sure the number of points is what we expect if we ask for more than we inserted
	results = idx.Nearby(origin, 10, AcceptAny)
	assertEqual(t, 8, len(results))
}

func TestKDTree_Nearby_NotMemphis(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex().Load(pts...)
	origin := NewCoordinates(-115, 45)

	results := idx.Nearby(origin, 3, func(pt Point) bool {
		return pt.(*NamedPoint).Name != "memphis"
	})

	assertEqual(t, 3, len(results))
	assertEqual(t, "woodinville", results[0].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[1].(*NamedPoint).Name)
	assertEqual(t, "anchorage", results[2].(*NamedPoint).Name)
}

func TestKDTree_Nearby_AntiMeridian(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex().Load(pts...)
	origin := NewCoordinates(-175, 60)

	results := idx.Nearby(origin, 3, AcceptAny)

	assertEqual(t,3, len(results))
	assertEqual(t, "eastrussia", results[0].(*NamedPoint).Name)
	assertEqual(t, "anchorage", results[1].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[2].(*NamedPoint).Name)
}

func TestKDTree_Nearby_NotEnough(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex().Load(pts...)
	origin := NewCoordinates(-175, 60)

	results := idx.Nearby(origin, 10, AcceptAny)

	assertEqual(t, 8, len(results))
}

func TestKDTree_Nearby_MultiNode(t *testing.T) {
	pts := namedPoints()
	// use a small NodeSize so our results must come from multiple nodes
	opts := KDTreeOptions{NodeSize: 2}
	idx := NewKDTreeIndex(opts).Load(pts...)
	origin := NewCoordinates(-175, -60)

	results := idx.Nearby(origin, 3, AcceptAny)

	assertEqual(t,3, len(results))
	assertEqual(t, "saopaulo", results[0].(*NamedPoint).Name)
	assertEqual(t, "tokyo", results[1].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[2].(*NamedPoint).Name)
}

func TestKDTree_Nearby_Picky(t *testing.T) {
	pts := globalPoints(100_000)
	idx := NewIndex().Load(pts...)

	// origin near north pole
	// only accept points in the southern hemisphere
	results := idx.Nearby(NewCoordinates(-1.23, 85), 1, func(pt Point) bool {
		return pt.Lat() < 0
	})
	assertEqual(t, 1, len(results))
	assertEqual(t, true, results[0].Lat() < 0)

	// origin near south pole
	// only accept points in the northern hemisphere
	results = idx.Nearby(NewCoordinates(-1.23, -85), 1, func(pt Point) bool {
		return pt.Lat() > 0
	})
	assertEqual(t, 1, len(results))
	assertEqual(t, true, results[0].Lat() > 0)
}

func TestKDTree_Nearby_Ranked(t *testing.T) {
	pts := []Point{
		&RankedPoint{
			Point: points["seattle"],
			Name:  "seattle-less-important",
			Rank:  1,
		},
		&RankedPoint{
			Point: points["seattle"],
			Name:  "seattle-more-important",
			Rank:  5,
		},
		&RankedPoint{
			Point: points["woodinville"],
			Name:  "woodinville-super-important",
			Rank:  5000,
		},
	}
	idx := NewIndex().Load(pts...)
	origin := NewCoordinates(-122, 47)

	results := idx.Nearby(origin, 3, AcceptAny)

	assertEqual(t, 3, len(results))
	// rank should make the more important one win if distance is equal
	assertEqual(t, "seattle-more-important", results[0].(*RankedPoint).Name)
	assertEqual(t, "seattle-less-important", results[1].(*RankedPoint).Name)
	// but distance still wins, regardless of rank
	assertEqual(t, "woodinville-super-important", results[2].(*RankedPoint).Name)
}

func TestKDTree_Nearby_Global(t *testing.T) {
	points := globalPoints(100_000)
	assertEqual(t, 100_000, len(points))
	idx := NewIndex().Load(points...)
	origin := NewCoordinates(-122, 47)

	results := idx.Nearby(origin, 5, AcceptAny)
	assertEqual(t, 5, len(results))
	for _, result := range results {
		distFromOrigin := distanceKm(origin, result)
		// we should find 5 points within 100 km
		assertEqual(t, true, distFromOrigin < 100)
	}
}

func TestKDTree_Nearby_Empty(t *testing.T) {
	idx := NewIndex()
	origin := NewCoordinates(-122, 47)
	results := idx.Nearby(origin, 5, AcceptAny)
	assertEqual(t, 0, len(results))
}

func TestKDTree_Load_Same(t *testing.T) {
	origin := NewCoordinates(-122, 47)
	points := namedPoints()

	idx := NewIndex().Load(points...)
	results := idx.Nearby(origin, 2, AcceptAny)
	assertEqual(t,2, len(results))
	assertEqual(t, "seattle", results[0].(*NamedPoint).Name)
	assertEqual(t, "woodinville", results[1].(*NamedPoint).Name)

	// load same points again
	idx.Load(points...)

	// same points, same results
	results = idx.Nearby(origin, 2, AcceptAny)
	assertEqual(t,2, len(results))
	assertEqual(t, "seattle", results[0].(*NamedPoint).Name)
	assertEqual(t, "woodinville", results[1].(*NamedPoint).Name)
}

func TestKDTree_Load_More(t *testing.T) {
	origin := NewCoordinates(-122, 47)
	points := namedPoints()

	idx := NewIndex().Load(points...)
	results := idx.Nearby(origin, 2, AcceptAny)
	assertEqual(t,2, len(results))
	assertEqual(t, "seattle", results[0].(*NamedPoint).Name)
	assertEqual(t, "woodinville", results[1].(*NamedPoint).Name)

	// duplicate the points and load again
	points = append(points, namedPoints()...)
	idx.Load(points...)

	// Should be 2X points now, there will be 2 Seattle points
	results = idx.Nearby(origin, 3, AcceptAny)
	assertEqual(t,3, len(results))
	assertEqual(t, "seattle", results[0].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[1].(*NamedPoint).Name)
	assertEqual(t, "woodinville", results[2].(*NamedPoint).Name)
}

func TestKDTree_Load_Less(t *testing.T) {
	origin := NewCoordinates(-122, 47)
	// start with duplicated points
	points := append(namedPoints(), namedPoints()...)

	idx := NewIndex().Load(points...)

	results := idx.Nearby(origin, 3, AcceptAny)
	// We start with 2X points, there will be only 2 Seattle points
	assertEqual(t,3, len(results))
	assertEqual(t, "seattle", results[0].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[1].(*NamedPoint).Name)
	assertEqual(t, "woodinville", results[2].(*NamedPoint).Name)

	// load again without duplicates
	idx.Load(namedPoints()...)

	// Should be only 1X points now, there will be only 1 Seattle
	results = idx.Nearby(origin, 2, AcceptAny)
	assertEqual(t,2, len(results))
	assertEqual(t, "seattle", results[0].(*NamedPoint).Name)
	assertEqual(t, "woodinville", results[1].(*NamedPoint).Name)
}
