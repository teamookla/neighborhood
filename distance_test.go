package neighborhood

import (
	"testing"
)

func TestHaversineMiles_Simple(t *testing.T) {
	pt1 := points["seattle"]
	pt2 := points["memphis"]
	dist := distance(pt1.Lon(), pt1.Lat(), pt2.Lon(), pt2.Lat())
	assertEqual(t, 3003, int(dist))
}

func TestHaversineMiles_AntiMeridian(t *testing.T) {
	// should go across the date line
	pt1 := points["anchorage"]
	pt2 := points["eastrussia"]
	dist := distance(pt1.Lon(), pt1.Lat(), pt2.Lon(), pt2.Lat())
	assertEqual(t, 1645, int(dist))
}

func TestHaversineMiles_Commute(t *testing.T) {
	pt1 := points["woodinville"]
	pt2 := points["seattle"]
	dist := distance(pt1.Lon(), pt1.Lat(), pt2.Lon(), pt2.Lat())
	assertEqual(t, 24, int(dist)) // not that far!
}
