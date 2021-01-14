package neighborhood

import (
	"math"
	"reflect"
	"testing"
)

var points = map[string]Point{
	"seattle":     NewCoordinates(-122.4, 47.6),
	"woodinville": NewCoordinates(-122.16, 47.75),
	"memphis":     NewCoordinates(-90.05, 35.15),
	"anchorage":   NewCoordinates(-150.0, 61.2),
	"tokyo":       NewCoordinates(139.67, 35.67),
	"eastrussia":  NewCoordinates(178.26, 63.06),
	"saopaulo":    NewCoordinates(-46.6, -23.5),
	"cairo":       NewCoordinates(31.2, 30.0),
}

type NamedPoint struct {
	Point
	Name string
}

type RankedPoint struct {
	Point
	Name string
	Rank float64
}

func (p RankedPoint) GetRank() float64 {
	return p.Rank
}

func namedPoint(name string) *NamedPoint {
	return &NamedPoint{
		Point: points[name],
		Name: name,
	}
}

func namedPoints() []Point {
	pts := make([]Point, 0, len(points))
	for name, _ := range points {
		pts = append(pts, namedPoint(name))
	}
	return pts
}

// globalPoints gets a list of n Points spread uniformly over the globe
func globalPoints(n int) []Point {
	var pts []Point
	maxStepsPerDimension := math.Ceil(math.Sqrt(float64(n)))
	latStep := 180.0 / maxStepsPerDimension
	lonStep := 360.0 / maxStepsPerDimension

	for lat := -90.0; lat <= 90; lat += latStep {
		for lon := -180.0; lon <= 180 && len(pts) < n; lon += lonStep {
			pts = append(pts, NewCoordinates(lon, lat))
		}
	}
	return pts
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
		t.FailNow()
	}
}

func assertNil(t *testing.T, a interface{}) {
	if a == nil || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil()) {
		return
	}
	t.Errorf("expected <nil>, got %v", a)
	t.FailNow()
}

