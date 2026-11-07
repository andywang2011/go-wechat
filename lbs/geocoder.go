package lbs

// Tencent Location Service WebService API endpoints and documentation links.
//
// Both reverse and forward geocoding use the same path; parameters differ.
const (
	// GeocoderV1URL shared GET endpoint for reverse and forward geocoding.
	GeocoderV1URL = "https://apis.map.qq.com/ws/geocoder/v1/"

	// ReverseGeocoderDocURL official doc for reverse geocoding (coordinates → address).
	ReverseGeocoderDocURL = "https://lbs.qq.com/service/webService/webServiceGuide/address/Gcoder"

	// ForwardGeocoderDocURL official doc for forward geocoding (address → coordinates).
	ForwardGeocoderDocURL = "https://lbs.qq.com/service/webService/webServiceGuide/address/Geocoder"
)

// GetPOI controls whether nearby POI list is returned (get_poi parameter).
const (
	GetPOIOff = "0" // omit (default)
	GetPOIOn  = "1" // include
)

// AddressPolicy forward-geocoding policy (policy parameter).
const (
	AddressPolicyStandard = 0 // standard: address must include a city name
	AddressPolicyRelaxed  = 1 // relaxed: city name may be omitted
)
