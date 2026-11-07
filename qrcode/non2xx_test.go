package qrcode

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-wechat/common"
)

func TestNon2xx_BeforeParse(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusInternalServerError} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"errcode":0,"errmsg":"ok"}`))
		}))
		api := NewClient(MapConfigProvider{"app1": {Host: srv.URL}})
		ctx := context.Background()

		var se *common.HTTPStatusError
		// Image endpoint.
		_, err := api.GetWxaCode(ctx, "app1", testAccessToken, GetWxaCodeReq{Path: "pages/index"})
		if !errors.As(err, &se) || se.StatusCode != status {
			t.Fatalf("GetWxaCode status %d: err = %v", status, err)
		}
		// JSON endpoint.
		_, err = api.GenShortLink(ctx, "app1", testAccessToken, GenShortLinkReq{PageURL: "pages/index"})
		if !errors.As(err, &se) || se.StatusCode != status {
			t.Fatalf("GenShortLink status %d: err = %v", status, err)
		}
		srv.Close()
	}
}

func TestImageTooLarge(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		buf := make([]byte, 1<<20)
		for i := 0; i < 11; i++ {
			_, _ = w.Write(buf)
		}
	}))
	defer srv.Close()
	api := NewClient(MapConfigProvider{"app1": {Host: srv.URL}})
	if _, err := api.GetWxaCode(context.Background(), "app1", testAccessToken, GetWxaCodeReq{Path: "p"}); err == nil {
		t.Fatal("expected size error")
	}
}
