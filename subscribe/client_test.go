package subscribe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendSubscript_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != SendSubscriptPath {
			t.Errorf("path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method: %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(SendSubscriptResp{RespErr: RespErr{Errcode: 0}})
	}))
	defer server.Close()

	cfg := MapConfigProvider{"app1": {Host: server.URL, AppID: "wx", Secret: "s"}}
	api := NewClient(cfg)
	resp, err := api.SendSubscript(context.Background(), "app1", "token",
		&SendSubscriptReq{Touser: "openid", TemplateID: "tpl", Data: map[string]SubscribeDataItem{"k": {Value: "v"}}})
	if err != nil {
		t.Fatalf("SendSubscript: %v", err)
	}
	if resp.Errcode != 0 {
		t.Fatalf("resp: %+v", resp)
	}
}

func TestClient_SendSubscript_WhitelistDenied(t *testing.T) {
	cfg := MapConfigProvider{"app1": {
		Host:               "https://example.com",
		AppID:              "wx",
		Secret:             "s",
		SubscribeWhitelist: []string{"allowed"},
	}}
	api := NewClient(cfg)
	_, err := api.SendSubscript(context.Background(), "app1", "token",
		&SendSubscriptReq{Touser: "blocked", TemplateID: "tpl"})
	if err == nil {
		t.Fatal("expected whitelist error")
	}
}

func TestMapConfigProvider_GetConfig(t *testing.T) {
	cfg := MapConfigProvider{"app1": {AppID: "wx", Secret: "s"}}
	got, err := cfg.GetConfig("app1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Host != DefaultHost {
		t.Errorf("host: %q", got.Host)
	}
}
