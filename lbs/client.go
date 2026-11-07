package lbs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/andywang2011/go-httplib"
	"github.com/andywang2011/go-wechat/common"
)

// API defines Tencent Location Service capabilities.
type API interface {
	LocationToAddress(ctx context.Context, req LocationToAddressReq) (*LocationToAddressResp, error)
	AddressToLocation(ctx context.Context, req AddressToLocationReq) (*AddressToLocationResp, error)
}

type geocoderRequest interface {
	apiURL() string
	queryParams() map[string]string
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client implements API.
type Client struct {
	cfg  Config
	doer doer
}

// NewClient creates an LBS client from Config.
func NewClient(cfg Config) (API, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Client{
		cfg:  cfg,
		doer: httplib.NewClient(httplib.WithTimeout(cfg.timeout())),
	}, nil
}

// NewClientWithDoer injects doer for tests.
func NewClientWithDoer(cfg Config, d doer) API {
	return &Client{cfg: cfg, doer: d}
}

func (c *Client) doRequest(ctx context.Context, lbs geocoderRequest, result interface{}) error {
	if c == nil || c.doer == nil {
		return fmt.Errorf("lbs: client not initialized")
	}
	if lbs == nil {
		return fmt.Errorf("lbs: request must not be nil")
	}
	data := lbs.queryParams()
	params := make(map[string]string, len(data)+2)
	for k, v := range data {
		params[k] = v
	}
	params["key"] = c.cfg.AppKey
	params["output"] = "json"

	rawURL := lbs.apiURL() + encodeQuery(params)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("lbs: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := common.ReadBody(resp, common.MaxResponseBytes)
	if err != nil {
		return fmt.Errorf("lbs: read body: %w", err)
	}
	if err := json.Unmarshal(b, result); err != nil {
		return fmt.Errorf("lbs: unmarshal: %w", err)
	}
	return nil
}

func encodeQuery(params map[string]string) string {
	v := url.Values{}
	for key, val := range params {
		v.Add(key, val)
	}
	return v.Encode()
}
