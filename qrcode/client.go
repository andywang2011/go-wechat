package qrcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/andywang2011/go-httplib"
	"github.com/andywang2011/go-wechat/common"
)

// API defines Mini Program QR code and link APIs.
type API interface {
	GetWxaCode(ctx context.Context, appkey, accessToken string, req GetWxaCodeReq) (*ImageResult, error)
	GetWxaCodeUnlimit(ctx context.Context, appkey, accessToken string, req GetWxaCodeUnlimitReq) (*ImageResult, error)
	CreateWxaQRCode(ctx context.Context, appkey, accessToken string, req CreateWxaQRCodeReq) (*ImageResult, error)
	QueryScheme(ctx context.Context, appkey, accessToken string, req QuerySchemeReq) (*QuerySchemeResp, error)
	GenerateScheme(ctx context.Context, appkey, accessToken string, req GenerateSchemeReq) (*GenerateSchemeResp, error)
	GenerateNFCScheme(ctx context.Context, appkey, accessToken string, req GenerateNFCSchemeReq) (*GenerateNFCSchemeResp, error)
	GenerateURLLink(ctx context.Context, appkey, accessToken string, req GenerateURLLinkReq) (*GenerateURLLinkResp, error)
	QueryURLLink(ctx context.Context, appkey, accessToken string, req QueryURLLinkReq) (*QueryURLLinkResp, error)
	GenShortLink(ctx context.Context, appkey, accessToken string, req GenShortLinkReq) (*GenShortLinkResp, error)
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client implements API.
type Client struct {
	config ConfigProvider
	client doer
}

// NewClient creates a qrcode API client.
func NewClient(config ConfigProvider) API {
	return &Client{
		config: config,
		client: httplib.NewClient(httplib.WithTimeout(10 * time.Second)),
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

// RespErr is the common WeChat API error envelope.
type RespErr struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

func (e RespErr) error() error {
	if e.Errcode == 0 {
		return nil
	}
	return fmt.Errorf("qrcode: errcode %d: %s", e.Errcode, e.Errmsg)
}

func (c *Client) postJSON(ctx context.Context, appkey, accessToken, path string, body, result interface{}) error {
	if accessToken == "" {
		return fmt.Errorf("qrcode: access_token is required")
	}
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return err
	}
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("qrcode: marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.host(cfg)+path+"?access_token="+url.QueryEscape(accessToken), bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doAndDecodeJSON(req, result)
}

func (c *Client) postImage(ctx context.Context, appkey, accessToken, path string, body interface{}) (*ImageResult, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("qrcode: access_token is required")
	}
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("qrcode: marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.host(cfg)+path+"?access_token="+url.QueryEscape(accessToken), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := common.ReadBody(resp, common.MaxImageResponseBytes)
	if err != nil {
		return nil, fmt.Errorf("qrcode: read body: %w", err)
	}
	ct := resp.Header.Get("Content-Type")
	if isJSONResponse(ct, data) {
		var errResp RespErr
		if err := json.Unmarshal(data, &errResp); err != nil {
			return nil, fmt.Errorf("qrcode: decode error response: %w", err)
		}
		return nil, errResp.error()
	}
	return &ImageResult{ContentType: ct, Data: data}, nil
}

func isJSONResponse(contentType string, data []byte) bool {
	if strings.Contains(contentType, "application/json") {
		return true
	}
	return len(data) > 0 && data[0] == '{'
}

func (c *Client) doAndDecodeJSON(req *http.Request, result interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := common.ReadBody(resp, common.MaxResponseBytes)
	if err != nil {
		return fmt.Errorf("qrcode: read body: %w", err)
	}
	if err := json.Unmarshal(b, result); err != nil {
		return fmt.Errorf("qrcode: unmarshal response: %w", err)
	}
	return nil
}
