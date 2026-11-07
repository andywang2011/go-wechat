package login

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testAccessToken = "test-access-token"

func TestSessionSignature(t *testing.T) {
	key := "test_session_key"
	got := SessionSignature(key)
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write(nil)
	want := hex.EncodeToString(mac.Sum(nil))
	if got != want {
		t.Errorf("SessionSignature = %q, want %q", got, want)
	}
}

func TestClient_Jscode2session_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != Jscode2sessionPath {
			t.Errorf("path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method: %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(Jscode2sessionResp{
			RespErr:    RespErr{Errcode: 0},
			Openid:     "openid_abc",
			Unionid:    "unionid_xyz",
			SessionKey: "session_key_123",
		})
	}))
	defer server.Close()

	cfg := MapConfigProvider{"app1": {Host: server.URL, AppID: "wx", Secret: "s"}}
	api := NewClient(cfg)
	resp, err := api.Jscode2session(context.Background(), "app1", "code_001")
	if err != nil {
		t.Fatalf("Jscode2session: %v", err)
	}
	if resp.Openid != "openid_abc" || resp.SessionKey != "session_key_123" {
		t.Errorf("resp: %+v", resp)
	}
}

func TestClient_CheckSession_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != CheckSessionPath {
			t.Errorf("path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("openid") != "openid_abc" {
			t.Errorf("openid: %s", r.URL.Query().Get("openid"))
		}
		if r.URL.Query().Get("sig_method") != sigMethodHMACSHA256 {
			t.Errorf("sig_method: %s", r.URL.Query().Get("sig_method"))
		}
		_ = json.NewEncoder(w).Encode(CheckSessionResp{RespErr: RespErr{Errcode: 0, Errmsg: "ok"}})
	}))
	defer server.Close()

	cfg := MapConfigProvider{"app1": {Host: server.URL, AppID: "wx", Secret: "s"}}
	api := NewClient(cfg)
	resp, err := api.CheckSession(context.Background(), "app1", testAccessToken, CheckSessionReq{
		OpenID:     "openid_abc",
		SessionKey: "session_key_123",
	})
	if err != nil {
		t.Fatalf("CheckSession: %v", err)
	}
	if resp.Errcode != 0 {
		t.Errorf("resp: %+v", resp)
	}
}

func TestClient_ResetSession_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != ResetSessionPath {
			t.Errorf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(ResetSessionResp{
			RespErr:    RespErr{Errcode: 0},
			Openid:     "openid_abc",
			SessionKey: "new_session_key",
		})
	}))
	defer server.Close()

	cfg := MapConfigProvider{"app1": {Host: server.URL, AppID: "wx", Secret: "s"}}
	api := NewClient(cfg)
	resp, err := api.ResetSession(context.Background(), "app1", testAccessToken, ResetSessionReq{
		OpenID:     "openid_abc",
		SessionKey: "session_key_123",
	})
	if err != nil {
		t.Fatalf("ResetSession: %v", err)
	}
	if resp.SessionKey != "new_session_key" {
		t.Errorf("resp: %+v", resp)
	}
}

func TestClient_CheckSession_EmptyToken(t *testing.T) {
	api := NewClient(MapConfigProvider{"app1": {AppID: "wx", Secret: "s"}})
	_, err := api.CheckSession(context.Background(), "app1", "", CheckSessionReq{OpenID: "o", SessionKey: "k"})
	if err == nil {
		t.Fatal("expected error for empty access_token")
	}
}

func TestClient_Jscode2session_ConfigError(t *testing.T) {
	api := NewClient(MapConfigProvider{})
	_, err := api.Jscode2session(context.Background(), "app1", "code")
	if err == nil {
		t.Fatal("expected error from config")
	}
}

func TestMapConfigProvider_GetConfig_DefaultHost(t *testing.T) {
	p := MapConfigProvider{"app1": {AppID: "wx_app1", Secret: "secret_app1"}}
	cfg, err := p.GetConfig("app1")
	if err != nil {
		t.Fatalf("GetConfig: %v", err)
	}
	if cfg.Host != DefaultHost {
		t.Errorf("expected default host %q, got %q", DefaultHost, cfg.Host)
	}
}
