package gis

func CoordsValid(lng, lat float64) bool {
	return ValidLng(lng) && ValidLat(lat)
}

func ValidLng(lng float64) bool {
	return lng >= -180.0 && lng <= 180.0
}

func ValidLat(lat float64) bool {
	return lat <= 90.0 && lat >= -90.0
}
