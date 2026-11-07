package lbs

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/andywang2011/go-wechat/common"
)

const locationOKBody = `{"status":0,"message":"ok","result":{"address":"addr","address_component":{"city":"city","district":"dist"},"ad_info":{"city_code":"010","adcode":"110000"}}}`

type mockDoer struct {
	body       string
	statusCode int
	err        error
}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
		Header:     make(http.Header),
	}, nil
}

func TestNewClient(t *testing.T) {
	cfg := Config{AppKey: "k", Timeout: 10 * time.Second}
	api, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if api == nil {
		t.Fatal("NewClient should return non-nil API")
	}
	_, ok := api.(*Client)
	if !ok {
		t.Error("NewClient should return *Client")
	}
}

func TestNewClient_ValidateError(t *testing.T) {
	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewClientWithDoer(t *testing.T) {
	cfg := Config{AppKey: "k"}
	m := &mockDoer{statusCode: 200, body: locationOKBody}
	c := NewClientWithDoer(cfg, m)
	if c == nil {
		t.Fatal("NewClientWithDoer should return non-nil client")
	}
}

func TestLocationToAddressReq_Validate(t *testing.T) {
	if err := (LocationToAddressReq{Lat: 39.9, Lng: 116.4}).Validate(); err != nil {
		t.Errorf("Validate: %v", err)
	}
	if err := (LocationToAddressReq{}).Validate(); err == nil {
		t.Error("expected error when lat/lng empty")
	}
	if err := (LocationToAddressReq{Lat: 1, Radius: 6000}).Validate(); err == nil {
		t.Error("expected error when radius > 5000")
	}
}

func TestLocationToAddressReq_queryParams(t *testing.T) {
	params := (LocationToAddressReq{
		Lat:        39.984154,
		Lng:        116.30749,
		Radius:     1000,
		GetPOI:     GetPOIOn,
		POIOptions: "policy=2",
	}).queryParams()
	if params["location"] != "39.984154,116.30749" {
		t.Errorf("location: %q", params["location"])
	}
	if params["radius"] != "1000" || params["get_poi"] != "1" || params["poi_options"] != "policy=2" {
		t.Errorf("params: %+v", params)
	}
	if (LocationToAddressReq{Lat: 1, Lng: 2}).apiURL() != GeocoderV1URL {
		t.Error("apiURL mismatch")
	}
}

func TestAddressToLocationReq_Validate(t *testing.T) {
	if err := (AddressToLocationReq{Address: "1 Main St"}).Validate(); err != nil {
		t.Errorf("Validate: %v", err)
	}
	if err := (AddressToLocationReq{}).Validate(); err == nil {
		t.Error("expected error when address empty")
	}
	if err := (AddressToLocationReq{Address: "x", Policy: 2}).Validate(); err == nil {
		t.Error("expected error for invalid policy")
	}
}

func TestAddressToLocationReq_queryParams(t *testing.T) {
	params := (AddressToLocationReq{
		Address: "1 Telegraph Hill Blvd, San Francisco, CA",
		Policy:  AddressPolicyRelaxed,
	}).queryParams()
	if params["address"] != "1 Telegraph Hill Blvd, San Francisco, CA" {
		t.Errorf("address: %q", params["address"])
	}
	if params["policy"] != "1" {
		t.Errorf("policy: %q", params["policy"])
	}
}

func TestClient_LocationToAddress_NilReceiver(t *testing.T) {
	var c *Client
	_, err := c.LocationToAddress(context.Background(), LocationToAddressReq{Lat: 1, Lng: 2})
	if err == nil {
		t.Error("expected error when receiver is nil")
	}
}

func TestClient_LocationToAddress_ValidateError(t *testing.T) {
	cfg := Config{AppKey: "k"}
	c := NewClientWithDoer(cfg, &mockDoer{}).(*Client)
	_, err := c.LocationToAddress(context.Background(), LocationToAddressReq{})
	if err == nil {
		t.Error("expected validation error")
	}
}

func TestClient_LocationToAddress_Success(t *testing.T) {
	cfg := Config{AppKey: "k"}
	c := NewClientWithDoer(cfg, &mockDoer{statusCode: 200, body: locationOKBody}).(*Client)
	resp, err := c.LocationToAddress(context.Background(), LocationToAddressReq{Lat: 39.9, Lng: 116.4})
	if err != nil {
		t.Fatalf("LocationToAddress err = %v", err)
	}
	if resp.Status != 0 || resp.Result.Address != "addr" || resp.Result.AdInfo.CityCode != "010" {
		t.Errorf("resp mismatch: %+v", resp)
	}
}

func TestAPI_LocationToAddress_Success(t *testing.T) {
	api := NewClientWithDoer(Config{AppKey: "k"}, &mockDoer{statusCode: 200, body: locationOKBody})
	resp, err := api.LocationToAddress(context.Background(), LocationToAddressReq{Lat: 39.9, Lng: 116.4})
	if err != nil {
		t.Fatalf("api.LocationToAddress err = %v", err)
	}
	if resp.Result.Address != "addr" {
		t.Errorf("resp.Result.Address = %q, want addr", resp.Result.Address)
	}
}

func TestClient_LocationToAddress_HTTPError(t *testing.T) {
	cfg := Config{AppKey: "k"}
	c := NewClientWithDoer(cfg, &mockDoer{statusCode: 401, body: "unauthorized"}).(*Client)
	_, err := c.LocationToAddress(context.Background(), LocationToAddressReq{Lat: 1, Lng: 2})
	if err == nil {
		t.Error("expected error for HTTP 401")
	}
}

func TestClient_LocationToAddress_DoError(t *testing.T) {
	cfg := Config{AppKey: "k"}
	c := NewClientWithDoer(cfg, &mockDoer{err: errors.New("do error")}).(*Client)
	_, err := c.LocationToAddress(context.Background(), LocationToAddressReq{Lat: 1, Lng: 2})
	if err == nil {
		t.Error("expected error when Do fails")
	}
}

func TestClient_AddressToLocation_Success(t *testing.T) {
	body := `{"status":0,"message":"ok","result":{"location":{"lng":116.4,"lat":39.9}}}`
	cfg := Config{AppKey: "k"}
	c := NewClientWithDoer(cfg, &mockDoer{statusCode: 200, body: body}).(*Client)
	resp, err := c.AddressToLocation(context.Background(), AddressToLocationReq{Address: "San Francisco"})
	if err != nil {
		t.Fatalf("AddressToLocation err = %v", err)
	}
	if resp.Status != 0 || resp.Result.Location.Lng != 116.4 || resp.Result.Location.Lat != 39.9 {
		t.Errorf("resp mismatch: %+v", resp)
	}
}

func TestClient_AddressToLocation_ValidateError(t *testing.T) {
	cfg := Config{AppKey: "k"}
	c := NewClientWithDoer(cfg, &mockDoer{}).(*Client)
	_, err := c.AddressToLocation(context.Background(), AddressToLocationReq{})
	if err == nil {
		t.Error("expected validation error")
	}
}

func TestDoRequest_NilLbs(t *testing.T) {
	cfg := Config{AppKey: "k"}
	c := NewClientWithDoer(cfg, &mockDoer{}).(*Client)
	var result LocationToAddressResp
	err := c.doRequest(context.Background(), nil, &result)
	if err == nil {
		t.Error("expected error for nil lbs")
	}
}

func TestConfig_Validate(t *testing.T) {
	if err := (Config{AppKey: "k"}).Validate(); err != nil {
		t.Errorf("Validate: %v", err)
	}
	if err := (Config{}).Validate(); err == nil {
		t.Error("expected error when AppKey empty")
	}
}

func TestConfig_timeout(t *testing.T) {
	if (Config{}).timeout() != defaultTimeout {
		t.Error("expected default timeout")
	}
	if (Config{Timeout: 5 * time.Second}).timeout() != 5*time.Second {
		t.Error("expected custom timeout")
	}
}

func TestFormatLocation(t *testing.T) {
	if got := formatLocation(39.9, 116.4); got != "39.9,116.4" {
		t.Errorf("formatLocation = %q", got)
	}
}

func TestClient_Non2xx_BeforeJSONParse(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusInternalServerError} {
		api := NewClientWithDoer(Config{AppKey: "k"}, &mockDoer{statusCode: status, body: locationOKBody})
		resp, err := api.LocationToAddress(context.Background(), LocationToAddressReq{Lat: 1, Lng: 2})
		var se *common.HTTPStatusError
		if !errors.As(err, &se) || se.StatusCode != status {
			t.Fatalf("status %d: err = %v", status, err)
		}
		if resp != nil {
			t.Fatalf("resp = %+v, want nil", resp)
		}
	}
}
