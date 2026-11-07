package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/andywang2011/go-httplib"
	"github.com/andywang2011/go-wechat/common"
)

// API defines Mini Program login server-side APIs.
type API interface {
	// Jscode2session exchanges js_code for openid and session_key (no access_token).
	Jscode2session(ctx context.Context, appkey, code string) (*Jscode2sessionResp, error)
	// CheckSession verifies a stored session_key is still valid (requires access_token).
	CheckSession(ctx context.Context, appkey, accessToken string, req CheckSessionReq) (*CheckSessionResp, error)
	// ResetSession resets session_key and returns a new one (requires access_token).
	ResetSession(ctx context.Context, appkey, accessToken string, req ResetSessionReq) (*ResetSessionResp, error)
}

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client implements API.
type Client struct {
	config ConfigProvider
	client doer
}

// NewClient creates a login API client.
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

func (c *Client) get(ctx context.Context, rawURL string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
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
		return fmt.Errorf("login: read body: %w", err)
	}
	if err := json.Unmarshal(b, result); err != nil {
		return fmt.Errorf("login: unmarshal response: %w", err)
	}
	return nil
}

// RespErr is the common WeChat API error envelope.
type RespErr struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

// CheckSessionReq identifies the user and proves possession of session_key via signature.
type CheckSessionReq struct {
	OpenID     string
	SessionKey string // used to compute signature; not sent to WeChat
	Signature  string // optional; computed from SessionKey when empty
	SigMethod  string // optional; defaults to hmac_sha256
}

// ResetSessionReq has the same signature fields as CheckSessionReq.
type ResetSessionReq struct {
	OpenID     string
	SessionKey string
	Signature  string
	SigMethod  string
}

func (r CheckSessionReq) signature() (sig, method string, err error) {
	return sessionSig(r.OpenID, r.SessionKey, r.Signature, r.SigMethod)
}

func (r ResetSessionReq) signature() (sig, method string, err error) {
	return sessionSig(r.OpenID, r.SessionKey, r.Signature, r.SigMethod)
}

func sessionSig(openid, sessionKey, signature, sigMethod string) (string, string, error) {
	if openid == "" {
		return "", "", fmt.Errorf("login: openid is required")
	}
	if sessionKey == "" && signature == "" {
		return "", "", fmt.Errorf("login: session_key or signature is required")
	}
	sig := signature
	if sig == "" {
		sig = SessionSignature(sessionKey)
	}
	method := sigMethod
	if method == "" {
		method = sigMethodHMACSHA256
	}
	return sig, method, nil
}

func (c *Client) sessionGet(
	ctx context.Context,
	cfg Config,
	accessToken, path string,
	q url.Values,
	result interface{},
) error {
	if accessToken == "" {
		return fmt.Errorf("login: access_token is required")
	}
	q.Set("access_token", accessToken)
	return c.get(ctx, c.host(cfg)+path+"?"+q.Encode(), result)
}
