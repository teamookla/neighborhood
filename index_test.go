package neighborhood

import (
	"fmt"
	"testing"
)

var cities = namedPoints()
func Example() {
	// Create a new index with all searchable `Points`. Indexes are static and immutable.
	idx := NewIndex(cities...)
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

	idx := NewIndex(thing)
	origin := NewCoordinates(-122, 47)
	results := idx.Nearby(origin, 5, AcceptAny)
	assertEqual(t, 1, len(results))
	assertEqual(t, "thing", results[0].(*Thing).name)
}

func TestKDBush_Nearby_Simple(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex(pts...)
	origin := NewCoordinates(-115, 45)

	results := idx.Nearby(origin, 3, AcceptAny)

	assertEqual(t, 3, len(results))
	assertEqual(t, "woodinville", results[0].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[1].(*NamedPoint).Name)
	assertEqual(t, "memphis", results[2].(*NamedPoint).Name)
}

func TestKDBush_Nearby_NotMemphis(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex(pts...)
	origin := NewCoordinates(-115, 45)

	results := idx.Nearby(origin, 3, func(pt Point) bool {
		return pt.(*NamedPoint).Name != "memphis"
	})

	assertEqual(t, 3, len(results))
	assertEqual(t, "woodinville", results[0].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[1].(*NamedPoint).Name)
	assertEqual(t, "anchorage", results[2].(*NamedPoint).Name)
}

func TestKDBush_Nearby_AntiMeridian(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex(pts...)
	origin := NewCoordinates(-175, 60)

	results := idx.Nearby(origin, 3, AcceptAny)

	assertEqual(t,3, len(results))
	assertEqual(t, "eastrussia", results[0].(*NamedPoint).Name)
	assertEqual(t, "anchorage", results[1].(*NamedPoint).Name)
	assertEqual(t, "seattle", results[2].(*NamedPoint).Name)
}

func TestKDBush_Nearby_NotEnough(t *testing.T) {
	pts := namedPoints()
	idx := NewIndex(pts...)
	origin := NewCoordinates(-175, 60)

	results := idx.Nearby(origin, 10, AcceptAny)

	assertEqual(t, 8, len(results))
}

func TestKDBush_Nearby_Ranked(t *testing.T) {
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
	idx := NewIndex(pts...)
	origin := NewCoordinates(-122, 47)

	results := idx.Nearby(origin, 3, AcceptAny)

	assertEqual(t, 3, len(results))
	// rank should make the more important one win if distance is equal
	assertEqual(t, "seattle-more-important", results[0].(*RankedPoint).Name)
	assertEqual(t, "seattle-less-important", results[1].(*RankedPoint).Name)
	// but distance still wins, regardless of rank
	assertEqual(t, "woodinville-super-important", results[2].(*RankedPoint).Name)
}

func TestKDBush_Nearby_Global(t *testing.T) {
	points := globalPoints()
	idx := NewIndex(points...)
	origin := NewCoordinates(-122, 47)

	results := idx.Nearby(origin, 5, AcceptAny)
	assertEqual(t, 5, len(results))
	assertEqual(t, -122.0, results[0].Lon())
	assertEqual(t, 47.0, results[0].Lat())
	assertEqual(t, -121.5, results[1].Lon())
	assertEqual(t, 47.0, results[1].Lat())
	assertEqual(t, -122.5, results[2].Lon())
	assertEqual(t, 47.0, results[2].Lat())
	assertEqual(t, -122.0, results[3].Lon())
	assertEqual(t, 47.5, results[3].Lat())
	assertEqual(t, -122.0, results[4].Lon())
	assertEqual(t, 46.5, results[4].Lat())
}

func TestKDBush_Nearby_Empty(t *testing.T) {
	idx := NewIndex()
	origin := NewCoordinates(-122, 47)
	results := idx.Nearby(origin, 5, AcceptAny)
	assertEqual(t, 0, len(results))
}