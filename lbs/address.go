package lbs

import (
	"context"
	"errors"
	"strconv"
)

// AddressToLocationReq forward geocoding request (text address → coordinates).
//
// API: GET GeocoderV1URL
// Doc: ForwardGeocoderDocURL
//
// Required: Address (province/city/district recommended; Client URL-encodes it).
// Key and output are appended by Client.
type AddressToLocationReq struct {
	// Address text to geocode; should include at least a city name for better accuracy.
	Address string
	// Policy parsing policy; see AddressPolicyStandard / AddressPolicyRelaxed.
	Policy int
	// Sig signature when WebService SN verification is enabled.
	Sig string
	// Callback JSONP callback function name.
	Callback string
}

// Validate checks required fields for forward geocoding.
func (r AddressToLocationReq) Validate() error {
	if r.Address == "" {
		return errors.New("lbs: address is required for forward geocoding")
	}
	if r.Policy != AddressPolicyStandard && r.Policy != AddressPolicyRelaxed {
		return errors.New("lbs: policy must be 0 (standard) or 1 (relaxed)")
	}
	return nil
}

func (r AddressToLocationReq) apiURL() string {
	return GeocoderV1URL
}

func (r AddressToLocationReq) queryParams() map[string]string {
	params := map[string]string{
		"address": r.Address,
	}
	if r.Policy != AddressPolicyStandard {
		params["policy"] = strconv.Itoa(r.Policy)
	}
	if r.Sig != "" {
		params["sig"] = r.Sig
	}
	if r.Callback != "" {
		params["callback"] = r.Callback
	}
	return params
}

type AddressToLocationResp struct {
	Status    int          `json:"status"`
	Message   string       `json:"message"`
	RequestId string       `json:"request_id"`
	Result    ResultDetail `json:"result"`
}

type ResultDetail struct {
	Location GeoDetail `json:"location"`
}

type GeoDetail struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

// AddressToLocation runs forward geocoding.
// Doc: ForwardGeocoderDocURL
func (c *Client) AddressToLocation(ctx context.Context, req AddressToLocationReq) (*AddressToLocationResp, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var result AddressToLocationResp
	if err := c.doRequest(ctx, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
