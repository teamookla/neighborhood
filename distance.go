package neighborhood

import "math"

const rad = math.Pi / 180.0

func haverSinDist(pt1 Point, lon2, lat2, cosLat1 float64) float64 {
	haverSinDLon := haverSin((pt1.Lon() - lon2) * rad)
	return haverSinDistPartial(haverSinDLon, cosLat1, pt1.Lat(), lat2)
}

// boxDist gets the lower bound for distance from a location to points inside a bounding box
func boxDist(pt Point, cosLat float64, node *kdTreeNode) float64 {
	// query point is between minimum and maximum longitudes
	if pt.Lon() >= node.MinLon && pt.Lon() <= node.MaxLon {
		if pt.Lat() < node.MinLat {
			return haverSin((pt.Lat() - node.MinLat) * rad)
		}
		if pt.Lat() > node.MaxLat {
			return haverSin((pt.Lat() - node.MaxLat) * rad)
		}
		return 0
	}

	// query point is west or east of the bounding box;
	// calculate the extremum for great circle distance from query point to the closest longitude;
	haverSinDLon := math.Min(haverSin((pt.Lon()-node.MinLon)*rad), haverSin((pt.Lon()-node.MaxLon)*rad))
	extremumLat := vertexLat(pt.Lat(), haverSinDLon)

	// if extremum is inside the box, return the distance to it
	if extremumLat > node.MinLat && extremumLat < node.MaxLat {
		return haverSinDistPartial(haverSinDLon, cosLat, pt.Lat(), extremumLat)
	}
	// otherwise return the distance to one of the bbox corners (whichever is closest)
	return math.Min(
		haverSinDistPartial(haverSinDLon, cosLat, pt.Lat(), node.MinLat),
		haverSinDistPartial(haverSinDLon, cosLat, pt.Lat(), node.MaxLat),
	)
}

func haverSin(theta float64) float64 {
	s := math.Sin(theta / 2)
	return math.Pow(s, 2)
}

func haverSinDistPartial(haverSinDLon, cosLat1, lat1, lat2 float64) float64 {
	return cosLat1*math.Cos(lat2*rad)*haverSinDLon + haverSin((lat1-lat2)*rad)
}

func vertexLat(lat, haverSinDLon float64) float64 {
	cosDLon := 1 - 2*haverSinDLon
	if cosDLon <= 0 {
		if lat > 0 {
			return 90
		}
		return -90
	}
	return math.Atan(math.Tan(lat*rad)/cosDLon) / rad
}
