package coupon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/andywang2011/go-httplib"
	"github.com/andywang2011/go-wechat/common"
)

// API defines WeChat coupon / card server APIs.
type API interface {
	// GetCode queries a user card code (requires access_token).
	GetCode(ctx context.Context, appkey, accessToken string, req *GetCodeReq) (*GetCodeResp, error)
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client implements API.
type Client struct {
	config ConfigProvider
	client doer
}

// NewClient creates a coupon API client.
func NewClient(config ConfigProvider) API {
	return &Client{
		config: config,
		client: httplib.NewClient(httplib.WithTimeout(3 * time.Second)),
	}
}

// NewClientWithDoer injects HTTP doer for tests.
func NewClientWithDoer(config ConfigProvider, d doer) API {
	return &Client{config: config, client: d}
}

func (c *Client) host(cfg Config) string {
	if cfg.Host != "" {
		return cfg.Host
	}
	return DefaultHost
}

func (c *Client) post(ctx context.Context, appkey, accessToken, path string, data, result interface{}) error {
	if accessToken == "" {
		return fmt.Errorf("coupon: access_token is required")
	}
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("coupon: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.host(cfg)+path+"?access_token="+url.QueryEscape(accessToken), bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doAndDecode(req, result)
}

func (c *Client) doAndDecode(req *http.Request, result interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := common.ReadBody(resp, common.MaxResponseBytes)
	if err != nil {
		return fmt.Errorf("coupon: read body: %w", err)
	}
	if err := json.Unmarshal(b, result); err != nil {
		return fmt.Errorf("coupon: unmarshal response: %w", err)
	}
	return nil
}

// RespErr is the common WeChat API error envelope.
type RespErr struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}
