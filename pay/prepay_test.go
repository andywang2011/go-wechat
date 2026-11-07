package pay

import (
	"context"
	"testing"
)

func basePrepay() PrepayBase {
	return PrepayBase{
		Description: "desc",
		OutTradeNo:  "ORDER123456",
		NotifyURL:   "https://example.com/notify",
		Amount:      Amount{Total: 100},
	}
}

func TestPrepayRequest_Validate_JSAPI(t *testing.T) {
	cfg := validTestConfig(t)
	req := PrepayRequest{
		Channel: ChannelJSAPI,
		Base:    basePrepay(),
		Payer:   &PayerInfo{OpenID: "openid"},
	}
	if err := req.Validate(cfg); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestPrepayRequest_Validate_JSAPI_NoOpenID(t *testing.T) {
	cfg := validTestConfig(t)
	req := PrepayRequest{Channel: ChannelJSAPI, Base: basePrepay()}
	if err := req.Validate(cfg); err == nil {
		t.Fatal("expected payer error")
	}
}

func TestPrepayRequest_Validate_H5_NoScene(t *testing.T) {
	cfg := validTestConfig(t)
	req := PrepayRequest{Channel: ChannelH5, Base: basePrepay()}
	if err := req.Validate(cfg); err == nil {
		t.Fatal("expected scene error")
	}
}

func TestPrepayRequest_Validate_InvalidNotifyURL(t *testing.T) {
	cfg := validTestConfig(t)
	b := basePrepay()
	b.NotifyURL = "http://insecure.com/notify"
	req := PrepayRequest{Channel: ChannelNative, Base: b}
	if err := req.Validate(cfg); err == nil {
		t.Fatal("expected https error")
	}
}

func TestPrepayRequest_Validate_Combine(t *testing.T) {
	cfg := validTestConfig(t)
	req := PrepayRequest{
		Channel: ChannelCombineJSAPI,
		Base:    basePrepay(),
		Payer:   &PayerInfo{OpenID: "openid"},
		Combine: &CombineInfo{
			CombineOutTradeNo: "COMBINE001",
			SubOrders: []SubOrder{{
				MchID: "sub_mch", Description: "sub", OutTradeNo: "SUB001", Amount: Amount{Total: 50},
			}},
		},
	}
	if err := req.Validate(cfg); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestPrepayRequest_Validate_OutTradeNoLength(t *testing.T) {
	cfg := validTestConfig(t)
	b := basePrepay()
	b.OutTradeNo = "short"
	req := PrepayRequest{Channel: ChannelAPP, Base: b}
	if err := req.Validate(cfg); err == nil {
		t.Fatal("expected length error")
	}
}

type mockGateway struct {
	lastChannel Channel
	prepayResp  *PrepayResponse
	err         error
}

func (m *mockGateway) prepay(ctx context.Context, req PrepayRequest, cfg Config) (*PrepayResponse, error) {
	m.lastChannel = req.Channel
	if m.err != nil {
		return nil, m.err
	}
	if m.prepayResp != nil {
		return m.prepayResp, nil
	}
	return &PrepayResponse{Channel: req.Channel, PrepayID: "prepay"}, nil
}

func (m *mockGateway) queryOrder(ctx context.Context, req QueryOrderRequest, cfg Config) (*Transaction, error) {
	return &Transaction{OutTradeNo: req.OutTradeNo}, nil
}

func (m *mockGateway) closeOrder(ctx context.Context, req CloseOrderRequest, cfg Config) error {
	return nil
}

func (m *mockGateway) refund(ctx context.Context, req RefundRequest, cfg Config) (*RefundResult, error) {
	return &RefundResult{OutRefundNo: req.OutRefundNo}, nil
}

func (m *mockGateway) queryRefund(ctx context.Context, req QueryRefundRequest, cfg Config) (*RefundResult, error) {
	return &RefundResult{OutRefundNo: req.OutRefundNo}, nil
}

func TestClient_Prepay_Dispatch(t *testing.T) {
	cfg := validTestConfig(t)
	mg := &mockGateway{}
	api := NewClientWithGateway(cfg, mg, nil)

	_, err := api.Prepay(context.Background(), PrepayRequest{
		Channel: ChannelNative,
		Base:    basePrepay(),
	})
	if err != nil {
		t.Fatalf("Prepay: %v", err)
	}
	if mg.lastChannel != ChannelNative {
		t.Errorf("channel = %q, want native", mg.lastChannel)
	}
}

func TestClient_QueryOrder(t *testing.T) {
	cfg := validTestConfig(t)
	api := NewClientWithGateway(cfg, &mockGateway{}, nil)
	tx, err := api.QueryOrder(context.Background(), QueryOrderRequest{OutTradeNo: "ORDER123456"})
	if err != nil {
		t.Fatalf("QueryOrder: %v", err)
	}
	if tx.OutTradeNo != "ORDER123456" {
		t.Errorf("OutTradeNo = %q", tx.OutTradeNo)
	}
}

func TestClient_Refund(t *testing.T) {
	cfg := validTestConfig(t)
	api := NewClientWithGateway(cfg, &mockGateway{}, nil)
	res, err := api.Refund(context.Background(), RefundRequest{
		OutTradeNo:  "ORDER123456",
		OutRefundNo: "REFUND001",
		Amount:      RefundAmount{Total: 100, Refund: 100},
	})
	if err != nil {
		t.Fatalf("Refund: %v", err)
	}
	if res.OutRefundNo != "REFUND001" {
		t.Errorf("OutRefundNo = %q", res.OutRefundNo)
	}
}
