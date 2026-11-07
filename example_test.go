package wechat_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/andywang2011/go-wechat/login"
)

func Example_loginJscode2session() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != login.Jscode2sessionPath {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(login.Jscode2sessionResp{
			RespErr:    login.RespErr{Errcode: 0},
			Openid:     "openid_demo",
			SessionKey: "session_key_demo",
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
	// Output: openid_demo
}
