package subscribe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/andywang2011/go-httplib"
	"github.com/andywang2011/go-wechat/common"
)

// API defines WeChat Mini Program subscribe message APIs.
type API interface {
	// SendSubscript sends a subscribe message (requires access_token).
	SendSubscript(ctx context.Context, appkey, accessToken string, req *SendSubscriptReq) (*SendSubscriptResp, error)
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client implements API.
type Client struct {
	config ConfigProvider
	client doer
}

// NewClient creates a subscribe message API client.
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
		return fmt.Errorf("subscribe: access_token is required")
	}
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("subscribe: marshal request: %w", err)
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
		return fmt.Errorf("subscribe: read body: %w", err)
	}
	if err := json.Unmarshal(b, result); err != nil {
		return fmt.Errorf("subscribe: unmarshal response: %w", err)
	}
	return nil
}

// RespErr is the common WeChat API error envelope.
type RespErr struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

func (c *Client) checkSubscribeWhitelist(cfg Config, openid string) error {
	if len(cfg.SubscribeWhitelist) == 0 {
		return nil
	}
	if !slices.Contains(cfg.SubscribeWhitelist, openid) {
		return fmt.Errorf("subscribe: openid %q not in subscribe whitelist", openid)
	}
	return nil
}

// SendSubscript sends a subscribe message after optional whitelist validation.
func (c *Client) SendSubscript(ctx context.Context, appkey, accessToken string, req *SendSubscriptReq) (*SendSubscriptResp, error) {
	if req == nil {
		return nil, fmt.Errorf("subscribe: request is required")
	}
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return nil, err
	}
	if err := c.checkSubscribeWhitelist(cfg, req.Touser); err != nil {
		return nil, err
	}
	var resp SendSubscriptResp
	if err := c.post(ctx, appkey, accessToken, SendSubscriptPath, req, &resp); err != nil {
		return nil, err
	}
	if resp.Errcode != 0 {
		return &resp, fmt.Errorf("subscribe: subscribe send error: %d %s", resp.Errcode, resp.Errmsg)
	}
	return &resp, nil
}
