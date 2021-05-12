package neighborhood

import (
	"math"
	"testing"
)

func TestHaversineMiles_Simple(t *testing.T) {
	pt1 := points["seattle"]
	pt2 := points["memphis"]
	dist := distanceKm(pt1, pt2)
	assertEqual(t, 3003, int(dist))
}

func TestHaversineMiles_AntiMeridian(t *testing.T) {
	// should go across the date line
	pt1 := points["anchorage"]
	pt2 := points["eastrussia"]
	dist := distanceKm(pt1, pt2)
	assertEqual(t, 1645, int(dist))
}

func TestHaversineMiles_Commute(t *testing.T) {
	pt1 := points["woodinville"]
	pt2 := points["seattle"]
	dist := distanceKm(pt1, pt2)
	assertEqual(t, 24, int(dist)) // not that far!
}


// distanceKm gets distance in kilometers for testing purposes
func distanceKm(pt1, pt2 Point) float64 {
	earthRadiusKm := 6371.0
	var h = haverSinDist(pt1, pt2.Lon(), pt2.Lat(), math.Cos(pt1.Lat()*rad))
	return 2 * earthRadiusKm * math.Asin(math.Sqrt(h))
}
