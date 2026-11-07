package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testAccessToken = "test-access-token"

func TestClient_GetUserPhoneNumber_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != GetUserPhoneNumberPath {
			t.Errorf("path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(GetUserPhoneNumberResp{
			RespErr: RespErr{Errcode: 0},
			PhoneInfo: PhoneInfo{
				PhoneNumber:     "13800138000",
				PurePhoneNumber: "13800138000",
				CountryCode:     "86",
			},
		})
	}))
	defer server.Close()

	cfg := MapConfigProvider{"app1": {Host: server.URL, AppID: "wx", Secret: "s"}}
	api := NewClient(cfg)
	resp, err := api.GetUserPhoneNumber(context.Background(), "app1", testAccessToken, "phone_code")
	if err != nil {
		t.Fatalf("GetUserPhoneNumber: %v", err)
	}
	if resp.PhoneInfo.PurePhoneNumber != "13800138000" {
		t.Errorf("resp: %+v", resp)
	}
}

func TestClient_GetUserPhoneNumber_EmptyToken(t *testing.T) {
	api := NewClient(MapConfigProvider{"app1": {AppID: "wx", Secret: "s"}})
	_, err := api.GetUserPhoneNumber(context.Background(), "app1", "", "code")
	if err == nil {
		t.Fatal("expected error for empty access_token")
	}
}
