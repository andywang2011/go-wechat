package login

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-wechat/common"
)

func non2xxServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		// Valid JSON that would parse as success if status were ignored.
		_, _ = w.Write([]byte(`{"openid":"should-not-parse","errcode":0}`))
	}))
	t.Cleanup(s.Close)
	return s
}

func assertStatusErr(t *testing.T, err error, want int) {
	t.Helper()
	var se *common.HTTPStatusError
	if !errors.As(err, &se) || se.StatusCode != want {
		t.Fatalf("err = %v, want HTTPStatusError %d", err, want)
	}
}

func TestNon2xx_BeforeJSONParse(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusInternalServerError, http.StatusBadGateway} {
		srv := non2xxServer(t, status)
		api := NewClient(MapConfigProvider{"app1": {Host: srv.URL, AppID: "wx", Secret: "s"}})
		ctx := context.Background()

		resp, err := api.Jscode2session(ctx, "app1", "code")
		assertStatusErr(t, err, status)
		if resp != nil && resp.Openid != "" {
			t.Errorf("body was parsed despite status %d: %+v", status, resp)
		}
		if _, err := api.CheckSession(ctx, "app1", testAccessToken, CheckSessionReq{OpenID: "o", SessionKey: "k"}); err == nil {
			t.Error("CheckSession expected error")
		} else {
			assertStatusErr(t, err, status)
		}
		if _, err := api.ResetSession(ctx, "app1", testAccessToken, ResetSessionReq{OpenID: "o", SessionKey: "k"}); err == nil {
			t.Error("ResetSession expected error")
		} else {
			assertStatusErr(t, err, status)
		}
	}
}
