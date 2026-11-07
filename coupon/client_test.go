package coupon

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testAccessToken = "test-access-token"

func TestClient_GetCode_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != GetCodePath {
			t.Errorf("path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method: %s", r.Method)
		}
		if got := r.URL.Query().Get("access_token"); got != testAccessToken {
			t.Errorf("access_token: %q", got)
		}
		_ = json.NewEncoder(w).Encode(GetCodeResp{
			RespErr: RespErr{Errcode: 0},
			CardID:  "card_123",
			OpenID:  "openid_demo",
		})
	}))
	defer server.Close()

	cfg := MapConfigProvider{"app1": {Host: server.URL, AppID: "wx", Secret: "s"}}
	api := NewClient(cfg)
	resp, err := api.GetCode(context.Background(), "app1", testAccessToken, &GetCodeReq{Code: "CODE1"})
	if err != nil {
		t.Fatalf("GetCode: %v", err)
	}
	if resp.OpenID != "openid_demo" || resp.CardID != "card_123" {
		t.Fatalf("resp: %+v", resp)
	}
}

func TestClient_GetCode_EmptyToken(t *testing.T) {
	api := NewClient(MapConfigProvider{"app1": {AppID: "wx", Secret: "s"}})
	_, err := api.GetCode(context.Background(), "app1", "", &GetCodeReq{Code: "c"})
	if err == nil {
		t.Fatal("expected error for empty access_token")
	}
}

func TestClient_GetCode_ConfigNotFound(t *testing.T) {
	api := NewClient(MapConfigProvider{})
	_, err := api.GetCode(context.Background(), "missing", testAccessToken, &GetCodeReq{Code: "c"})
	if err == nil {
		t.Fatal("expected config error")
	}
}

func TestClient_GetCode_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	cfg := MapConfigProvider{"app1": {Host: server.URL, AppID: "wx", Secret: "s"}}
	api := NewClient(cfg)
	_, err := api.GetCode(context.Background(), "app1", testAccessToken, &GetCodeReq{Code: "c"})
	if err == nil {
		t.Fatal("expected unmarshal error on invalid JSON body")
	}
}

func TestMapConfigProvider_GetConfig_DefaultHost(t *testing.T) {
	cfg := MapConfigProvider{"app1": {AppID: "wx", Secret: "s"}}
	got, err := cfg.GetConfig("app1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Host != DefaultHost {
		t.Errorf("host: %q", got.Host)
	}
}

func TestMapConfigProvider_GetConfig_EmptyAppkey(t *testing.T) {
	cfg := MapConfigProvider{"app1": {AppID: "wx"}}
	if _, err := cfg.GetConfig(""); err == nil {
		t.Fatal("expected error for empty appkey")
	}
}

type errDoer struct{ err error }

func (e errDoer) Do(_ *http.Request) (*http.Response, error) {
	return nil, e.err
}

func TestClient_GetCode_TransportError(t *testing.T) {
	api := NewClientWithDoer(MapConfigProvider{"app1": {AppID: "wx", Secret: "s"}}, errDoer{err: context.Canceled})
	_, err := api.GetCode(context.Background(), "app1", testAccessToken, &GetCodeReq{Code: "c"})
	if err == nil {
		t.Fatal("expected transport error")
	}
}
