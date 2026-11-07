package pay

import (
	"context"
	"net/http"
)

// API is the unified WeChat Pay API v3 client surface.
type API interface {
	Prepay(ctx context.Context, req PrepayRequest) (*PrepayResponse, error)
	QueryOrder(ctx context.Context, req QueryOrderRequest) (*Transaction, error)
	CloseOrder(ctx context.Context, req CloseOrderRequest) error
	Refund(ctx context.Context, req RefundRequest) (*RefundResult, error)
	QueryRefund(ctx context.Context, req QueryRefundRequest) (*RefundResult, error)
	ParseNotify(ctx context.Context, r *http.Request, dst any) (*NotifySummary, error)
}
