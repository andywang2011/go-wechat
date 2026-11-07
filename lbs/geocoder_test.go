package lbs

import "testing"

func TestGeocoderConstants(t *testing.T) {
	if GeocoderV1URL != "https://apis.map.qq.com/ws/geocoder/v1/" {
		t.Errorf("GeocoderV1URL = %q", GeocoderV1URL)
	}
	if ReverseGeocoderDocURL == "" || ForwardGeocoderDocURL == "" {
		t.Error("doc URLs must not be empty")
	}
}
