package lbs

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

const (
	KeySecondCountLimit = 120
	KeyDayCountLimit    = 121
)

// LocationToAddressReq reverse geocoding request (coordinates → text address).
//
// API: GET GeocoderV1URL
// Doc: ReverseGeocoderDocURL
//
// Required: Lat, Lng (GCJ-02; query location=lat,lng).
// Key and output are appended by Client.
type LocationToAddressReq struct {
	// Lat latitude (GCJ-02).
	Lat float64
	// Lng longitude (GCJ-02).
	Lng float64

	// Radius administrative snap radius in meters, 0–5000, default 0; increase offshore when no district.
	Radius int
	// GetPOI whether to return nearby POIs; see GetPOIOff / GetPOIOn.
	GetPOI string
	// POIOptions nearby POI options, e.g. address_format=short;radius=5000;policy=2.
	POIOptions string
	// Sig signature when WebService SN verification is enabled; see official FAQ.
	Sig string
	// Callback JSONP callback function name.
	Callback string
}

// Validate checks required fields for reverse geocoding.
func (r LocationToAddressReq) Validate() error {
	if r.Lat == 0 && r.Lng == 0 {
		return errors.New("lbs: lat and lng are required for reverse geocoding")
	}
	if r.Radius < 0 || r.Radius > 5000 {
		return errors.New("lbs: radius must be between 0 and 5000")
	}
	return nil
}

func (r LocationToAddressReq) apiURL() string {
	return GeocoderV1URL
}

func (r LocationToAddressReq) queryParams() map[string]string {
	params := map[string]string{
		"location": formatLocation(r.Lat, r.Lng),
	}
	if r.Radius > 0 {
		params["radius"] = strconv.Itoa(r.Radius)
	}
	if r.GetPOI != "" {
		params["get_poi"] = r.GetPOI
	}
	if r.POIOptions != "" {
		params["poi_options"] = r.POIOptions
	}
	if r.Sig != "" {
		params["sig"] = r.Sig
	}
	if r.Callback != "" {
		params["callback"] = r.Callback
	}
	return params
}

func formatLocation(lat, lng float64) string {
	return fmt.Sprintf("%v,%v", lat, lng)
}

type LocationToAddressResp struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Result    struct {
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
		Address          string `json:"address"`
		AddressComponent struct {
			Nation       string `json:"nation"`
			Province     string `json:"province"`
			City         string `json:"city"`
			District     string `json:"district"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
		} `json:"address_component"`
		AdInfo struct {
			NationCode    string `json:"nation_code"`
			Adcode        string `json:"adcode"`
			PhoneAreaCode string `json:"phone_area_code"`
			CityCode      string `json:"city_code"`
			Name          string `json:"name"`
			Location      struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			Nation   string `json:"nation"`
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
		} `json:"ad_info"`
	} `json:"result"`
}

// LocationToAddress runs reverse geocoding.
// Doc: ReverseGeocoderDocURL
func (c *Client) LocationToAddress(ctx context.Context, req LocationToAddressReq) (*LocationToAddressResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var result LocationToAddressResp
	if err := c.doRequest(ctx, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
