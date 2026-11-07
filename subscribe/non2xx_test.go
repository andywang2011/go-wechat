package subscribe

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andywang2011/go-wechat/common"
)

func TestNon2xx_BeforeJSONParse(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusServiceUnavailable} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{"errcode":0}`))
		}))
		api := NewClient(MapConfigProvider{"app1": {Host: srv.URL, AppID: "wx", Secret: "s"}})
		resp, err := api.SendSubscript(context.Background(), "app1", "token",
			&SendSubscriptReq{Touser: "o", TemplateID: "t"})
		srv.Close()
		var se *common.HTTPStatusError
		if !errors.As(err, &se) || se.StatusCode != status {
			t.Fatalf("status %d: err = %v", status, err)
		}
		if resp != nil {
			t.Fatalf("resp = %+v, want nil", resp)
		}
	}
}
