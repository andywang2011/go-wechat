package pay

import (
	"fmt"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

func resolveAppID(base PrepayBase, cfg Config) string {
	if base.AppID != "" {
		return base.AppID
	}
	return cfg.DefaultAppID
}

func parseTimeExpire(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("pay: time_expire must be RFC3339: %w", err)
	}
	return &t, nil
}

func toSDKAmount(a Amount) (total int64, currency *string) {
	total = a.Total
	cur := a.Currency
	if cur == "" {
		cur = "CNY"
	}
	return total, core.String(cur)
}

func toJSAPIAmount(a Amount) *jsapi.Amount {
	total, cur := toSDKAmount(a)
	return &jsapi.Amount{Total: &total, Currency: cur}
}

func toNativeAmount(a Amount) *native.Amount {
	total, cur := toSDKAmount(a)
	return &native.Amount{Total: &total, Currency: cur}
}

func toH5Amount(a Amount) *h5.Amount {
	total, cur := toSDKAmount(a)
	return &h5.Amount{Total: &total, Currency: cur}
}

func toAPPAmount(a Amount) *app.Amount {
	total, cur := toSDKAmount(a)
	return &app.Amount{Total: &total, Currency: cur}
}

func buildJSAPIPrepay(req PrepayRequest, cfg Config) (jsapi.PrepayRequest, error) {
	appid := resolveAppID(req.Base, cfg)
	if appid == "" {
		return jsapi.PrepayRequest{}, fmt.Errorf("pay: app_id is required")
	}
	expire, err := parseTimeExpire(req.Base.TimeExpire)
	if err != nil {
		return jsapi.PrepayRequest{}, err
	}
	out := jsapi.PrepayRequest{
		Appid:       core.String(appid),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(req.Base.Description),
		OutTradeNo:  core.String(req.Base.OutTradeNo),
		NotifyUrl:   core.String(req.Base.NotifyURL),
		Amount:      toJSAPIAmount(req.Base.Amount),
	}
	if expire != nil {
		out.TimeExpire = expire
	}
	if req.Base.Attach != "" {
		out.Attach = core.String(req.Base.Attach)
	}
	if req.Base.GoodsTag != "" {
		out.GoodsTag = core.String(req.Base.GoodsTag)
	}
	if req.Payer != nil && req.Payer.OpenID != "" {
		out.Payer = &jsapi.Payer{Openid: core.String(req.Payer.OpenID)}
	}
	return out, nil
}

func buildNativePrepay(req PrepayRequest, cfg Config) (native.PrepayRequest, error) {
	appid := resolveAppID(req.Base, cfg)
	if appid == "" {
		return native.PrepayRequest{}, fmt.Errorf("pay: app_id is required")
	}
	expire, err := parseTimeExpire(req.Base.TimeExpire)
	if err != nil {
		return native.PrepayRequest{}, err
	}
	out := native.PrepayRequest{
		Appid:       core.String(appid),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(req.Base.Description),
		OutTradeNo:  core.String(req.Base.OutTradeNo),
		NotifyUrl:   core.String(req.Base.NotifyURL),
		Amount:      toNativeAmount(req.Base.Amount),
	}
	if expire != nil {
		out.TimeExpire = expire
	}
	if req.Base.Attach != "" {
		out.Attach = core.String(req.Base.Attach)
	}
	if req.Base.GoodsTag != "" {
		out.GoodsTag = core.String(req.Base.GoodsTag)
	}
	return out, nil
}

func buildH5Prepay(req PrepayRequest, cfg Config) (h5.PrepayRequest, error) {
	appid := resolveAppID(req.Base, cfg)
	if appid == "" {
		return h5.PrepayRequest{}, fmt.Errorf("pay: app_id is required")
	}
	expire, err := parseTimeExpire(req.Base.TimeExpire)
	if err != nil {
		return h5.PrepayRequest{}, err
	}
	out := h5.PrepayRequest{
		Appid:       core.String(appid),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(req.Base.Description),
		OutTradeNo:  core.String(req.Base.OutTradeNo),
		NotifyUrl:   core.String(req.Base.NotifyURL),
		Amount:      toH5Amount(req.Base.Amount),
	}
	if expire != nil {
		out.TimeExpire = expire
	}
	if req.Base.Attach != "" {
		out.Attach = core.String(req.Base.Attach)
	}
	if req.Base.GoodsTag != "" {
		out.GoodsTag = core.String(req.Base.GoodsTag)
	}
	if req.Scene != nil {
		out.SceneInfo = &h5.SceneInfo{
			PayerClientIp: core.String("127.0.0.1"),
			H5Info: &h5.H5Info{
				Type: core.String(req.Scene.Type),
			},
		}
		if req.Scene.AppName != "" {
			out.SceneInfo.H5Info.AppName = core.String(req.Scene.AppName)
		}
		if req.Scene.AppURL != "" {
			out.SceneInfo.H5Info.AppUrl = core.String(req.Scene.AppURL)
		}
		if req.Scene.BundleID != "" {
			out.SceneInfo.H5Info.BundleId = core.String(req.Scene.BundleID)
		}
		if req.Scene.PackageName != "" {
			out.SceneInfo.H5Info.PackageName = core.String(req.Scene.PackageName)
		}
	}
	return out, nil
}

func buildAPPPrepay(req PrepayRequest, cfg Config) (app.PrepayRequest, error) {
	appid := resolveAppID(req.Base, cfg)
	if appid == "" {
		return app.PrepayRequest{}, fmt.Errorf("pay: app_id is required")
	}
	expire, err := parseTimeExpire(req.Base.TimeExpire)
	if err != nil {
		return app.PrepayRequest{}, err
	}
	out := app.PrepayRequest{
		Appid:       core.String(appid),
		Mchid:       core.String(cfg.MchID),
		Description: core.String(req.Base.Description),
		OutTradeNo:  core.String(req.Base.OutTradeNo),
		NotifyUrl:   core.String(req.Base.NotifyURL),
		Amount:      toAPPAmount(req.Base.Amount),
	}
	if expire != nil {
		out.TimeExpire = expire
	}
	if req.Base.Attach != "" {
		out.Attach = core.String(req.Base.Attach)
	}
	if req.Base.GoodsTag != "" {
		out.GoodsTag = core.String(req.Base.GoodsTag)
	}
	return out, nil
}

func fromJSAPIInvoke(r *jsapi.PrepayWithRequestPaymentResponse) *JSAPIInvokeParams {
	if r == nil {
		return nil
	}
	return &JSAPIInvokeParams{
		AppID:     deref(r.Appid),
		TimeStamp: deref(r.TimeStamp),
		NonceStr:  deref(r.NonceStr),
		Package:   deref(r.Package),
		SignType:  deref(r.SignType),
		PaySign:   deref(r.PaySign),
	}
}

func fromTransaction(t *payments.Transaction) *Transaction {
	if t == nil {
		return nil
	}
	out := &Transaction{
		AppID:          deref(t.Appid),
		MchID:          deref(t.Mchid),
		OutTradeNo:     deref(t.OutTradeNo),
		TransactionID:  deref(t.TransactionId),
		TradeType:      deref(t.TradeType),
		TradeState:     deref(t.TradeState),
		TradeStateDesc: deref(t.TradeStateDesc),
		Attach:         deref(t.Attach),
		SuccessTime:    deref(t.SuccessTime),
		BankType:       deref(t.BankType),
	}
	if t.Amount != nil && t.Amount.Total != nil {
		out.Amount = &Amount{Total: *t.Amount.Total, Currency: deref(t.Amount.Currency)}
	}
	if t.Payer != nil {
		out.PayerOpenID = deref(t.Payer.Openid)
	}
	return out
}

func toRefundCreate(req RefundRequest) refunddomestic.CreateRequest {
	cur := req.Amount.Currency
	if cur == "" {
		cur = "CNY"
	}
	out := refunddomestic.CreateRequest{
		OutRefundNo: core.String(req.OutRefundNo),
		Amount: &refunddomestic.AmountReq{
			Refund:   core.Int64(req.Amount.Refund),
			Total:    core.Int64(req.Amount.Total),
			Currency: core.String(cur),
		},
	}
	if req.OutTradeNo != "" {
		out.OutTradeNo = core.String(req.OutTradeNo)
	}
	if req.Reason != "" {
		out.Reason = core.String(req.Reason)
	}
	if req.NotifyURL != "" {
		out.NotifyUrl = core.String(req.NotifyURL)
	}
	return out
}

func fromRefund(r *refunddomestic.Refund) *RefundResult {
	if r == nil {
		return nil
	}
	out := &RefundResult{
		RefundID:            deref(r.RefundId),
		OutRefundNo:         deref(r.OutRefundNo),
		TransactionID:       deref(r.TransactionId),
		OutTradeNo:          deref(r.OutTradeNo),
		UserReceivedAccount: deref(r.UserReceivedAccount),
		Status:              refundStatusString(r.Status),
	}
	if r.SuccessTime != nil {
		out.SuccessTime = r.SuccessTime.Format(time.RFC3339)
	}
	if r.CreateTime != nil {
		out.CreateTime = r.CreateTime.Format(time.RFC3339)
	}
	if r.Amount != nil {
		out.Amount = &RefundAmount{
			Total:    derefInt(r.Amount.Total),
			Refund:   derefInt(r.Amount.Refund),
			Currency: deref(r.Amount.Currency),
		}
	}
	return out
}

func refundStatusString(s *refunddomestic.Status) string {
	if s == nil {
		return ""
	}
	return string(*s)
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
