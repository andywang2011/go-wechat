# go-wechat

Tencent WeChat product SDKs for Go (multi-package module).

Subpackages: `common` (access_token fetch), `coupon` (reserved), `lbs` (geocoding), `login` (mini program login), `miniprogram` (shared config via `subscribe`), `pay` (WeChat Pay API v3), `qrcode` (codes and links), `subscribe` (subscribe messages), `user` (user info).

## Install

```bash
go get github.com/andywang2011/go-wechat@v0.1.0
```

## Quick start

```go
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
```

## Examples

```bash
cd examples
go run .
```

## License

MIT — see [LICENSE](LICENSE).
