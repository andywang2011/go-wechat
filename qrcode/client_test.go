package qrcode

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testAccessToken = "test-access-token"

func TestClient_GetWxaCodeUnlimit_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathGetWxaCodeUnlimit {
			t.Errorf("path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte{0x89, 0x50, 0x4e, 0x47})
	}))
	defer server.Close()

	api := NewClient(MapConfigProvider{"app1": {Host: server.URL}})
	img, err := api.GetWxaCodeUnlimit(context.Background(), "app1", testAccessToken, GetWxaCodeUnlimitReq{
		Scene: "id=1",
		Page:  "pages/index",
	})
	if err != nil {
		t.Fatalf("GetWxaCodeUnlimit: %v", err)
	}
	if len(img.Data) < 4 {
		t.Errorf("expected image bytes, got %d", len(img.Data))
	}
}

func TestClient_GetWxaCode_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(RespErr{Errcode: 40001, Errmsg: "invalid credential"})
	}))
	defer server.Close()

	api := NewClient(MapConfigProvider{"app1": {Host: server.URL}})
	_, err := api.GetWxaCode(context.Background(), "app1", testAccessToken, GetWxaCodeReq{Path: "pages/index"})
	if err == nil {
		t.Fatal("expected API error")
	}
}

func TestClient_GenerateScheme_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathGenerateScheme {
			t.Errorf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(GenerateSchemeResp{
			RespErr:  RespErr{Errcode: 0},
			OpenLink: "weixin://dl/business/?t=xxx",
		})
	}))
	defer server.Close()

	api := NewClient(MapConfigProvider{"app1": {Host: server.URL}})
	resp, err := api.GenerateScheme(context.Background(), "app1", testAccessToken, GenerateSchemeReq{
		JumpWxa: &JumpWxa{Path: "pages/index", Query: "a=1"},
	})
	if err != nil {
		t.Fatalf("GenerateScheme: %v", err)
	}
	if resp.OpenLink == "" {
		t.Errorf("resp: %+v", resp)
	}
}

func TestClient_GenShortLink_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathGenShortLink {
			t.Errorf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(GenShortLinkResp{
			RespErr: RespErr{Errcode: 0},
			Link:    "#miniapp://demo",
		})
	}))
	defer server.Close()

	api := NewClient(MapConfigProvider{"app1": {Host: server.URL}})
	resp, err := api.GenShortLink(context.Background(), "app1", testAccessToken, GenShortLinkReq{
		PageURL: "pages/index",
	})
	if err != nil {
		t.Fatalf("GenShortLink: %v", err)
	}
	if resp.Link == "" {
		t.Errorf("resp: %+v", resp)
	}
}

func TestDocURL(t *testing.T) {
	if DocURL(PathGetWxaCode) == "" {
		t.Error("doc url should be set")
	}
}
