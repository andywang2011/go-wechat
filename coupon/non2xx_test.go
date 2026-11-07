package coupon

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-wechat/common"
)

func TestNon2xx_BeforeJSONParse(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusInternalServerError} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"errcode":0}`))
		}))
		api := NewClient(MapConfigProvider{"app1": {Host: srv.URL, AppID: "wx", Secret: "s"}})
		_, err := api.GetCode(context.Background(), "app1", testAccessToken, &GetCodeReq{Code: "c"})
		srv.Close()
		var se *common.HTTPStatusError
		if !errors.As(err, &se) || se.StatusCode != status {
			t.Fatalf("status %d: err = %v", status, err)
		}
	}
}
