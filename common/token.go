package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// AccessTokenPath GET /cgi-bin/token — shared by mini program, official account, etc.
	// Doc: https://developers.weixin.qq.com/miniprogram/dev/server/API/mp-access-token/api_getaccesstoken.html
	AccessTokenPath = "/cgi-bin/token"
)

// AccessTokenProvider returns access_token by appkey.
// Implementations (memory, Redis) live in the application layer.
type AccessTokenProvider interface {
	GetAccessToken(ctx context.Context, appkey string) (string, error)
}

// RespErr is the common WeChat API error envelope.
type RespErr struct {
	// Errcode is 0 on success.
	Errcode int `json:"errcode,omitempty"`
	// Errmsg is the error message when Errcode is non-zero.
	Errmsg string `json:"errmsg,omitempty"`
}

// AccessTokenResult is the response from FetchAccessToken.
type AccessTokenResult struct {
	RespErr
	// AccessToken is the credential used for subsequent Open API calls.
	AccessToken string `json:"access_token"`
	// ExpiresIn is the token lifetime in seconds.
	ExpiresIn int `json:"expires_in"`
}

// FetchAccessToken performs a stateless HTTP call to obtain access_token.
func FetchAccessToken(ctx context.Context, cred AppCredential, client *http.Client) (*AccessTokenResult, error) {
	if client == nil {
		return nil, fmt.Errorf("common: http client is required")
	}
	host := cred.HostOrDefault()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s%s?grant_type=client_credential&appid=%s&secret=%s",
			host, AccessTokenPath, cred.AppID, cred.Secret), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ReadBody(resp, MaxResponseBytes)
	if err != nil {
		return nil, fmt.Errorf("common: read body: %w", err)
	}
	var body AccessTokenResult
	if err := json.Unmarshal(b, &body); err != nil {
		return nil, fmt.Errorf("common: unmarshal: %w", err)
	}
	if body.Errcode != 0 {
		return nil, fmt.Errorf("common: token error: %d %s", body.Errcode, body.Errmsg)
	}
	if body.AccessToken == "" {
		return nil, fmt.Errorf("common: access_token empty")
	}
	return &body, nil
}

// CacheTTL derives cache duration from expires_in with a safety margin.
func CacheTTL(expiresIn int) int {
	ttl := expiresIn - 120
	if ttl <= 0 {
		ttl = 1800
	}
	if ttl > 86400 {
		ttl = 86400
	}
	return ttl
}
