package login

import (
	"context"
	"fmt"
	"net/url"
)

const (
	// Jscode2sessionPath GET /sns/jscode2session — Mini Program login credential check
	// Doc: https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	Jscode2sessionPath = "/sns/jscode2session"

	// CheckSessionPath GET /wxa/checksession — verify login session
	// Doc: https://developers.weixin.qq.com/miniprogram/dev/server/API/user-login/api_checksessionkey.html
	CheckSessionPath = "/wxa/checksession"

	// ResetSessionPath GET /wxa/resetusersessionkey — reset login session
	// Doc: https://developers.weixin.qq.com/miniprogram/dev/server/API/user-login/api_resetusersessionkey.html
	ResetSessionPath = "/wxa/resetusersessionkey"

	SystemError      = -1
	InvalidCodeError = 40029
	CodeUsedError    = 40163
)

// Jscode2sessionResp is the code2Session response.
type Jscode2sessionResp struct {
	RespErr
	Openid     string `json:"openid,omitempty"`
	Unionid    string `json:"unionid,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
}

// CheckSessionResp is the checkSession response (usually errcode/errmsg only).
type CheckSessionResp struct {
	RespErr
}

// ResetSessionResp is the resetUserSessionKey response.
type ResetSessionResp struct {
	RespErr
	Openid     string `json:"openid,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
}

// Jscode2session exchanges js_code for session_key and openid (no access_token required).
func (c *Client) Jscode2session(ctx context.Context, appkey, jscode string) (*Jscode2sessionResp, error) {
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("appid", cfg.AppID)
	q.Set("secret", cfg.Secret)
	q.Set("js_code", jscode)
	q.Set("grant_type", "authorization_code")

	var resp Jscode2sessionResp
	err = c.get(ctx, c.host(cfg)+Jscode2sessionPath+"?"+q.Encode(), &resp)
	return &resp, err
}

// Jscode2sessionURL builds the request URL for documentation/testing.
func Jscode2sessionURL(host, appID, secret, jsCode string) string {
	if host == "" {
		host = DefaultHost
	}
	return fmt.Sprintf("%s%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		host, Jscode2sessionPath, url.QueryEscape(appID), url.QueryEscape(secret), url.QueryEscape(jsCode))
}

// CheckSession verifies the server-stored session_key via signature (requires access_token).
func (c *Client) CheckSession(ctx context.Context, appkey, accessToken string, req CheckSessionReq) (*CheckSessionResp, error) {
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return nil, err
	}
	sig, method, err := req.signature()
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("openid", req.OpenID)
	q.Set("signature", sig)
	q.Set("sig_method", method)

	var resp CheckSessionResp
	err = c.sessionGet(ctx, cfg, accessToken, CheckSessionPath, q, &resp)
	return &resp, err
}

// ResetSession resets session_key and returns a new session_key (requires access_token).
func (c *Client) ResetSession(ctx context.Context, appkey, accessToken string, req ResetSessionReq) (*ResetSessionResp, error) {
	cfg, err := c.config.GetConfig(appkey)
	if err != nil {
		return nil, err
	}
	sig, method, err := req.signature()
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Set("openid", req.OpenID)
	q.Set("signature", sig)
	q.Set("sig_method", method)

	var resp ResetSessionResp
	err = c.sessionGet(ctx, cfg, accessToken, ResetSessionPath, q, &resp)
	return &resp, err
}
