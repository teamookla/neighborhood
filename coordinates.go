package neighborhood

// Coordinates are a simple Point with latitude and longitude
type Coordinates struct {
	lon float64
	lat float64
}

// NewCoordinates creates a new simple Point with specified latitude and longitude
func NewCoordinates(lon, lat float64) Point {
	return Coordinates{lon, lat}
}

// implement Point interface

// Lon gets the coordinate longitude
func (c Coordinates) Lon() float64 { return c.lon }

// Lat gets the coordinate latitude
func (c Coordinates) Lat() float64 { return c.lat }
