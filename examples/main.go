package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/andywang2011/go-wechat/login"
)

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(login.Jscode2sessionResp{
			RespErr:    login.RespErr{Errcode: 0},
			Openid:     "example-openid",
			SessionKey: "example-session-key",
		})
	}))
	defer srv.Close()

	client := login.NewClient(login.MapConfigProvider{
		"app1": {Host: srv.URL, AppID: "wx_appid", Secret: "app_secret"},
	})
	resp, err := client.Jscode2session(context.Background(), "app1", "js_code")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("openid:", resp.Openid)
}
