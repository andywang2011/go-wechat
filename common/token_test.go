package common

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAccessToken_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != AccessTokenPath {
			t.Errorf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(AccessTokenResult{
			RespErr:     RespErr{Errcode: 0},
			AccessToken: "mock_token_123",
			ExpiresIn:   7200,
		})
	}))
	defer server.Close()

	result, err := FetchAccessToken(context.Background(), AppCredential{
		Host: server.URL, AppID: "wx", Secret: "secret",
	}, server.Client())
	if err != nil {
		t.Fatalf("FetchAccessToken: %v", err)
	}
	if result.AccessToken != "mock_token_123" {
		t.Errorf("token: %s", result.AccessToken)
	}
}

func TestFetchAccessToken_NilClient(t *testing.T) {
	_, err := FetchAccessToken(context.Background(), AppCredential{AppID: "wx", Secret: "s"}, nil)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestFetchAccessToken_Errcode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(AccessTokenResult{
			RespErr: RespErr{Errcode: 40013, Errmsg: "invalid appid"},
		})
	}))
	defer server.Close()

	_, err := FetchAccessToken(context.Background(), AppCredential{
		Host: server.URL, AppID: "wx", Secret: "secret",
	}, server.Client())
	if err == nil {
		t.Fatal("expected error when errcode != 0")
	}
}

func TestFetchAccessToken_EmptyToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(AccessTokenResult{
			RespErr:     RespErr{Errcode: 0},
			AccessToken: "",
			ExpiresIn:   7200,
		})
	}))
	defer server.Close()

	_, err := FetchAccessToken(context.Background(), AppCredential{
		Host: server.URL, AppID: "wx", Secret: "secret",
	}, server.Client())
	if err == nil {
		t.Fatal("expected error when access_token empty")
	}
}

func TestCacheTTL(t *testing.T) {
	if CacheTTL(7200) != 7080 {
		t.Errorf("CacheTTL(7200) = %d", CacheTTL(7200))
	}
	if CacheTTL(60) != 1800 {
		t.Errorf("CacheTTL(60) = %d", CacheTTL(60))
	}
	if CacheTTL(100000) != 86400 {
		t.Errorf("CacheTTL(100000) = %d", CacheTTL(100000))
	}
}

func TestAppCredential_HostOrDefault(t *testing.T) {
	if (AppCredential{}).HostOrDefault() != DefaultAPIHost {
		t.Error("expected default host")
	}
	if (AppCredential{Host: "https://custom"}).HostOrDefault() != "https://custom" {
		t.Error("expected custom host")
	}
}
