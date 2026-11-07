package login_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/andywang2011/go-wechat/login"
)

// Jscode2session exchanges wx.login code for openid and session_key (no access_token).
func ExampleClient_Jscode2session() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(login.Jscode2sessionResp{
			RespErr:    login.RespErr{Errcode: 0},
			Openid:     "openid_abc",
			SessionKey: "session_key_123",
		})
	}))
	defer srv.Close()

	client := login.NewClient(login.MapConfigProvider{
		"app1": {Host: srv.URL, AppID: "wx_appid", Secret: "app_secret"},
	})
	resp, err := client.Jscode2session(context.Background(), "app1", "js_code_from_wx_login")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Openid)
	// Output: openid_abc
}

// CheckSession verifies a stored session_key is still valid.
func ExampleClient_CheckSession() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(login.CheckSessionResp{
			RespErr: login.RespErr{Errcode: 0, Errmsg: "ok"},
		})
	}))
	defer srv.Close()

	client := login.NewClient(login.MapConfigProvider{
		"app1": {Host: srv.URL, AppID: "wx", Secret: "s"},
	})
	_, err := client.CheckSession(context.Background(), "app1", "access-token", login.CheckSessionReq{
		OpenID:     "user_openid",
		SessionKey: "stored_session_key",
	})
	if err != nil {
		panic(err)
	}
}

func ExampleSessionSignature() {
	fmt.Println(login.CheckSessionPath)
	// Output: /wxa/checksession
}
