package pay

import (
	"context"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"

	"github.com/andywang2011/go-wechat/common"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/consts"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type paymentGateway interface {
	prepay(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error)
	queryOrder(ctx context.Context, req QueryOrderRequest, cfg Config) (*Transaction, error)
	closeOrder(ctx context.Context, req CloseOrderRequest, cfg Config) error
	refund(ctx context.Context, req RefundRequest, cfg Config) (*RefundResult, error)
	queryRefund(ctx context.Context, req QueryRefundRequest, cfg Config) (*RefundResult, error)
}

type sdkServices struct {
	jsapi  *jsapi.JsapiApiService
	native *native.NativeApiService
	h5     *h5.H5ApiService
	app    *app.AppApiService
	refund *refunddomestic.RefundsApiService
	http   *core.Client
}

func newSDKServices(client *core.Client) sdkServices {
	return sdkServices{
		jsapi:  &jsapi.JsapiApiService{Client: client},
		native: &native.NativeApiService{Client: client},
		h5:     &h5.H5ApiService{Client: client},
		app:    &app.AppApiService{Client: client},
		refund: &refunddomestic.RefundsApiService{Client: client},
		http:   client,
	}
}

type sdkGateway struct {
	svc sdkServices
}

func newSDKGateway(client *core.Client) *sdkGateway {
	return &sdkGateway{svc: newSDKServices(client)}
}

func (g *sdkGateway) prepay(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error) {
	switch req.Channel {
	case ChannelJSAPI:
		return g.prepayJSAPI(ctx, req, cfg)
	case ChannelNative:
		return g.prepayNative(ctx, req, cfg)
	case ChannelH5:
		return g.prepayH5(ctx, req, cfg)
	case ChannelAPP:
		return g.prepayAPP(ctx, req, cfg)
	case ChannelCombineJSAPI, ChannelCombineNative, ChannelCombineH5, ChannelCombineAPP:
		return g.prepayCombine(ctx, req, cfg)
	default:
		return nil, fmt.Errorf("pay: unsupported channel %q", req.Channel)
	}
}

func (g *sdkGateway) prepayJSAPI(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error) {
	sdkReq, err := buildJSAPIPrepay(req, cfg)
	if err != nil {
		return nil, err
	}
	resp, _, err := g.svc.jsapi.PrepayWithRequestPayment(ctx, sdkReq)
	if err != nil {
		return nil, err
	}
	return &PrepayResponse{
		Channel:  ChannelJSAPI,
		PrepayID: deref(resp.PrepayId),
		JSAPI:    fromJSAPIInvoke(resp),
	}, nil
}

func (g *sdkGateway) prepayNative(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error) {
	sdkReq, err := buildNativePrepay(req, cfg)
	if err != nil {
		return nil, err
	}
	resp, _, err := g.svc.native.Prepay(ctx, sdkReq)
	if err != nil {
		return nil, err
	}
	return &PrepayResponse{
		Channel: ChannelNative,
		Native:  &NativeResult{CodeURL: deref(resp.CodeUrl)},
	}, nil
}

func (g *sdkGateway) prepayH5(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error) {
	sdkReq, err := buildH5Prepay(req, cfg)
	if err != nil {
		return nil, err
	}
	resp, _, err := g.svc.h5.Prepay(ctx, sdkReq)
	if err != nil {
		return nil, err
	}
	return &PrepayResponse{
		Channel: ChannelH5,
		H5:      &H5Result{H5URL: deref(resp.H5Url)},
	}, nil
}

func (g *sdkGateway) prepayAPP(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error) {
	sdkReq, err := buildAPPPrepay(req, cfg)
	if err != nil {
		return nil, err
	}
	resp, _, err := g.svc.app.Prepay(ctx, sdkReq)
	if err != nil {
		return nil, err
	}
	pid := deref(resp.PrepayId)
	return &PrepayResponse{
		Channel:  ChannelAPP,
		PrepayID: pid,
		APP:      &APPResult{PrepayID: pid},
	}, nil
}

type combinePrepayBody struct {
	CombineAppid      string            `json:"combine_appid"`
	CombineMchid      string            `json:"combine_mchid"`
	CombineOutTradeNo string            `json:"combine_out_trade_no"`
	NotifyURL         string            `json:"notify_url"`
	TimeExpire        string            `json:"time_expire,omitempty"`
	SubOrders         []combineSubOrder `json:"sub_orders"`
	CombinePayerInfo  *combinePayerInfo `json:"combine_payer_info,omitempty"`
	SceneInfo         *combineSceneInfo `json:"scene_info,omitempty"`
}

type combineSubOrder struct {
	Mchid       string        `json:"mchid"`
	Attach      string        `json:"attach,omitempty"`
	Description string        `json:"description"`
	OutTradeNo  string        `json:"out_trade_no"`
	Amount      combineAmount `json:"amount"`
}

type combineAmount struct {
	Total    int64  `json:"total"`
	Currency string `json:"currency,omitempty"`
}

type combinePayerInfo struct {
	OpenID string `json:"openid"`
}

type combineSceneInfo struct {
	H5Info *combineH5Info `json:"h5_info,omitempty"`
}

type combineH5Info struct {
	Type        string `json:"type"`
	AppName     string `json:"app_name,omitempty"`
	AppURL      string `json:"app_url,omitempty"`
	BundleID    string `json:"bundle_id,omitempty"`
	PackageName string `json:"package_name,omitempty"`
}

type combinePrepayResp struct {
	PrepayID string `json:"prepay_id"`
	CodeURL  string `json:"code_url"`
	H5URL    string `json:"h5_url"`
}

func buildCombineBody(req PrepayRequest, cfg Config) (combinePrepayBody, error) {
	appid := resolveAppID(req.Base, cfg)
	if appid == "" {
		return combinePrepayBody{}, fmt.Errorf("pay: app_id is required")
	}
	if req.Combine == nil {
		return combinePrepayBody{}, fmt.Errorf("pay: combine info is required")
	}
	body := combinePrepayBody{
		CombineAppid:      appid,
		CombineMchid:      cfg.MchID,
		CombineOutTradeNo: req.Combine.CombineOutTradeNo,
		NotifyURL:         req.Base.NotifyURL,
		TimeExpire:        req.Base.TimeExpire,
	}
	for _, sub := range req.Combine.SubOrders {
		cur := sub.Amount.Currency
		if cur == "" {
			cur = "CNY"
		}
		body.SubOrders = append(body.SubOrders, combineSubOrder{
			Mchid:       sub.MchID,
			Attach:      sub.Attach,
			Description: sub.Description,
			OutTradeNo:  sub.OutTradeNo,
			Amount:      combineAmount{Total: sub.Amount.Total, Currency: cur},
		})
	}
	switch req.Channel {
	case ChannelCombineJSAPI:
		if req.Payer == nil || req.Payer.OpenID == "" {
			return combinePrepayBody{}, fmt.Errorf("pay: payer openid is required for combine jsapi")
		}
		body.CombinePayerInfo = &combinePayerInfo{OpenID: req.Payer.OpenID}
	case ChannelCombineH5:
		if req.Scene == nil || req.Scene.Type == "" {
			return combinePrepayBody{}, fmt.Errorf("pay: h5 scene is required for combine h5")
		}
		body.SceneInfo = &combineSceneInfo{H5Info: &combineH5Info{
			Type:        req.Scene.Type,
			AppName:     req.Scene.AppName,
			AppURL:      req.Scene.AppURL,
			BundleID:    req.Scene.BundleID,
			PackageName: req.Scene.PackageName,
		}}
	}
	return body, nil
}

func (g *sdkGateway) prepayCombine(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error) {
	body, err := buildCombineBody(req, cfg)
	if err != nil {
		return nil, err
	}
	var resp combinePrepayResp
	if err := g.postJSON(ctx, http.MethodPost, req.Channel.APIPath(), body, &resp); err != nil {
		return nil, err
	}
	out := &PrepayResponse{
		Channel:  req.Channel,
		PrepayID: resp.PrepayID,
		Combine: &CombinePrepayResult{
			PrepayID: resp.PrepayID,
			CodeURL:  resp.CodeURL,
			H5URL:    resp.H5URL,
		},
	}
	switch req.Channel {
	case ChannelCombineNative:
		out.Native = &NativeResult{CodeURL: resp.CodeURL}
	case ChannelCombineH5:
		out.H5 = &H5Result{H5URL: resp.H5URL}
	case ChannelCombineAPP:
		out.APP = &APPResult{PrepayID: resp.PrepayID}
	}
	return out, nil
}

func (g *sdkGateway) postJSON(ctx context.Context, method, path string, body, resp interface{}) error {
	result, err := g.svc.http.Request(ctx, method, consts.WechatPayAPIServer+path, nil, nil, body, "application/json")
	if err != nil {
		return err
	}
	if result.Response == nil {
		return fmt.Errorf("pay: empty response for %s", path)
	}
	defer result.Response.Body.Close()
	if result.Response.StatusCode < 200 || result.Response.StatusCode >= 300 {
		return fmt.Errorf("pay: %s returned status %d", path, result.Response.StatusCode)
	}
	if resp != nil {
		limitResponseBody(result.Response)
		if err := core.UnMarshalResponse(result.Response, resp); err != nil {
			return fmt.Errorf("pay: decode response: %w", err)
		}
	}
	return nil
}

func (g *sdkGateway) queryOrder(ctx context.Context, req QueryOrderRequest, cfg Config) (*Transaction, error) {
	if req.Kind == OrderKindCombine {
		return g.queryCombineOrder(ctx, req, cfg)
	}
	sdkReq := jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(cfg.MchID),
	}
	resp, _, err := g.svc.jsapi.QueryOrderByOutTradeNo(ctx, sdkReq)
	if err != nil {
		return nil, err
	}
	return fromTransaction(resp), nil
}

func (g *sdkGateway) queryCombineOrder(ctx context.Context, req QueryOrderRequest, cfg Config) (*Transaction, error) {
	path := fmt.Sprintf("/v3/combine-transactions/out-trade-no/%s?mchid=%s",
		neturl.PathEscape(req.OutTradeNo), neturl.QueryEscape(cfg.MchID))
	var raw map[string]interface{}
	if err := g.getJSON(ctx, path, &raw); err != nil {
		return nil, err
	}
	return combineMapToTransaction(raw), nil
}

func (g *sdkGateway) getJSON(ctx context.Context, path string, resp interface{}) error {
	result, err := g.svc.http.Request(ctx, http.MethodGet, consts.WechatPayAPIServer+path, nil, nil, nil, "")
	if err != nil {
		return err
	}
	defer result.Response.Body.Close()
	if result.Response.StatusCode < 200 || result.Response.StatusCode >= 300 {
		return fmt.Errorf("pay: GET %s status %d", path, result.Response.StatusCode)
	}
	limitResponseBody(result.Response)
	return core.UnMarshalResponse(result.Response, resp)
}

// limitResponseBody caps how much of the response body decoders may read.
func limitResponseBody(r *http.Response) {
	r.Body = limitedBody{Reader: io.LimitReader(r.Body, common.MaxResponseBytes), Closer: r.Body}
}

type limitedBody struct {
	io.Reader
	io.Closer
}

func combineMapToTransaction(m map[string]interface{}) *Transaction {
	if m == nil {
		return nil
	}
	out := &Transaction{}
	if v, ok := m["combine_out_trade_no"].(string); ok {
		out.OutTradeNo = v
	}
	if v, ok := m["combine_mchid"].(string); ok {
		out.MchID = v
	}
	if v, ok := m["combine_appid"].(string); ok {
		out.AppID = v
	}
	return out
}

func (g *sdkGateway) closeOrder(ctx context.Context, req CloseOrderRequest, cfg Config) error {
	if req.Kind == OrderKindCombine {
		path := fmt.Sprintf("/v3/combine-transactions/out-trade-no/%s/close", neturl.PathEscape(req.OutTradeNo))
		body := map[string]string{"combine_mchid": cfg.MchID}
		return g.postJSON(ctx, http.MethodPost, path, body, nil)
	}
	_, err := g.svc.jsapi.CloseOrder(ctx, jsapi.CloseOrderRequest{
		OutTradeNo: core.String(req.OutTradeNo),
		Mchid:      core.String(cfg.MchID),
	})
	return err
}

func (g *sdkGateway) refund(ctx context.Context, req RefundRequest, cfg Config) (*RefundResult, error) {
	resp, _, err := g.svc.refund.Create(ctx, toRefundCreate(req))
	if err != nil {
		return nil, err
	}
	return fromRefund(resp), nil
}

func (g *sdkGateway) queryRefund(ctx context.Context, req QueryRefundRequest, cfg Config) (*RefundResult, error) {
	resp, _, err := g.svc.refund.QueryByOutRefundNo(ctx, refunddomestic.QueryByOutRefundNoRequest{
		OutRefundNo: core.String(req.OutRefundNo),
	})
	if err != nil {
		return nil, err
	}
	return fromRefund(resp), nil
}
