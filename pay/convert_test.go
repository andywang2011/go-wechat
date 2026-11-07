package pay

import (
	"testing"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

func TestBuildJSAPIPrepay(t *testing.T) {
	cfg := validTestConfig(t)
	req := PrepayRequest{
		Channel: ChannelJSAPI,
		Base: PrepayBase{
			Description: "d",
			OutTradeNo:  "ORDER123456",
			NotifyURL:   "https://example.com/n",
			Amount:      Amount{Total: 1},
			TimeExpire:  "2015-05-20T13:29:35+08:00",
			GoodsTag:    "tag",
		},
		Payer: &PayerInfo{OpenID: "oid"},
	}
	sdk, err := buildJSAPIPrepay(req, cfg)
	if err != nil {
		t.Fatalf("buildJSAPIPrepay: %v", err)
	}
	if deref(sdk.Appid) != cfg.DefaultAppID {
		t.Errorf("appid = %q", deref(sdk.Appid))
	}
	if deref(sdk.Payer.Openid) != "oid" {
		t.Errorf("openid = %q", deref(sdk.Payer.Openid))
	}
}

func TestFromTransaction(t *testing.T) {
	total := int64(100)
	tx := fromTransaction(&payments.Transaction{
		OutTradeNo: core.String("OUT1"),
		Amount:     &payments.TransactionAmount{Total: &total, Currency: core.String("CNY")},
	})
	if tx.OutTradeNo != "OUT1" || tx.Amount.Total != 100 {
		t.Fatalf("unexpected tx: %+v", tx)
	}
}
