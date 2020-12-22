package neighborhood

import "math"

const earthRadius = 6371.0
const rad = math.Pi / 180.0

func distance(lng1, lat1, lng2, lat2 float64) float64 {
	var h = haverSinDist(lng1, lat1, lng2, lat2, math.Cos(lat1*rad))
	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}

// boxDist gets the lower bound for distance from a location to points inside a bounding box
func boxDist(lon, lat, cosLat float64, node *kdTreeNode) float64 {
	minLon := node.MinLon
	maxLon := node.MaxLon
	minLat := node.MinLat
	maxLat := node.MaxLat

	// query point is between minimum and maximum longitudes
	if lon >= minLon && lon <= maxLon {
		if lat < minLat {
			return haverSin((lat - minLat) * rad)
		}
		if lat > maxLat {
			return haverSin((lat - maxLat) * rad)
		}
		return 0
	}

	// query point is west or east of the bounding box;
	// calculate the extremum for great circle distance from query point to the closest longitude;
	haverSinDLng := math.Min(haverSin((lon-minLon)*rad), haverSin((lon-maxLon)*rad))
	extremumLat := vertexLat(lat, haverSinDLng)

	// if extremum is inside the box, return the distance to it
	if extremumLat > minLat && extremumLat < maxLat {
		return haverSinDistPartial(haverSinDLng, cosLat, lat, extremumLat)
	}
	// otherwise return the distance to one of the bbox corners (whichever is closest)
	return math.Min(
		haverSinDistPartial(haverSinDLng, cosLat, lat, minLat),
		haverSinDistPartial(haverSinDLng, cosLat, lat, maxLat),
	)
}

func haverSin(theta float64) float64 {
	s := math.Sin(theta / 2)
	return math.Pow(s, 2)
}

func haverSinDistPartial(haverSinDLng, cosLat1, lat1, lat2 float64) float64 {
	return cosLat1*math.Cos(lat2*rad)*haverSinDLng + haverSin((lat1-lat2)*rad)
}

func haverSinDist(lng1, lat1, lng2, lat2, cosLat1 float64) float64 {
	haverSinDLng := haverSin((lng1 - lng2) * rad)
	return haverSinDistPartial(haverSinDLng, cosLat1, lat1, lat2)
}

func vertexLat(lat, haverSinDLng float64) float64 {
	cosDLng := 1 - 2*haverSinDLng
	if cosDLng <= 0 {
		if lat > 0 {
			return 90
		}
		return -90
	}
	return math.Atan(math.Tan(lat*rad)/cosDLng) / rad
}
