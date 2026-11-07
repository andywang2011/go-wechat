package pay

import "time"

// OrderKind distinguishes normal vs combine order lifecycle APIs.
type OrderKind int

const (
	// OrderKindNormal is a single-merchant order.
	OrderKindNormal OrderKind = iota
	// OrderKindCombine is a combine (multi-sub-order) payment.
	OrderKindCombine
)

// Amount is an order amount in the smallest currency unit (CNY: fen).
type Amount struct {
	// Total is the order total in the smallest currency unit.
	Total int64
	// Currency is the ISO 4217 code; defaults to CNY when empty.
	Currency string
}

// PrepayBase holds fields common to all prepay channels.
type PrepayBase struct {
	// Description is the order title shown to the payer.
	Description string
	// OutTradeNo is the merchant-side unique order number.
	OutTradeNo string
	// NotifyURL receives payment result notifications.
	NotifyURL string
	// Amount is the order amount block.
	Amount Amount
	// Attach is optional merchant metadata echoed in callbacks.
	Attach string
	// TimeExpire is optional RFC3339 expiry (e.g. 2015-05-20T13:29:35+08:00).
	TimeExpire string
	// AppID overrides Config.DefaultAppID when set.
	AppID string
	// GoodsTag is an optional order discount tag.
	GoodsTag string
}

// PayerInfo is required for JSAPI / combine JSAPI.
type PayerInfo struct {
	// OpenID is the payer's WeChat OpenID under the paying app.
	OpenID string
}

// H5SceneInfo is required for H5 and combine H5.
type H5SceneInfo struct {
	// Type is the scene type (iOS, Android, Wap).
	Type string
	// AppName is the app or site name shown in the payment page.
	AppName string
	// AppURL is the official site URL for Wap type.
	AppURL string
	// BundleID is the iOS bundle identifier for iOS type.
	BundleID string
	// PackageName is the Android package name for Android type.
	PackageName string
}

// SubOrder is one sub-order in a combine payment.
type SubOrder struct {
	// MchID is the sub-merchant WeChat Pay merchant ID.
	MchID string
	// Attach is optional merchant metadata for this sub-order.
	Attach string
	// Description is the sub-order title.
	Description string
	// OutTradeNo is the merchant-side unique sub-order number.
	OutTradeNo string
	// Amount is the sub-order amount block.
	Amount Amount
}

// CombineInfo is required for combine-* channels.
type CombineInfo struct {
	// CombineOutTradeNo is the merchant-side combine order number.
	CombineOutTradeNo string
	// SubOrders lists each sub-order in the combine payment.
	SubOrders []SubOrder
}

// PrepayRequest is the unified prepay input; Channel selects the API.
type PrepayRequest struct {
	// Channel selects JSAPI, Native, H5, APP, or combine variants.
	Channel Channel
	// Base holds common prepay fields.
	Base PrepayBase
	// Payer is required for JSAPI and combine JSAPI channels.
	Payer *PayerInfo
	// Scene is required for H5 and combine H5 channels.
	Scene *H5SceneInfo
	// Combine is required for combine payment channels.
	Combine *CombineInfo
}

// JSAPIInvokeParams are returned to the front end for JSAPI / mini program pay.
type JSAPIInvokeParams struct {
	// AppID is the WeChat app ID.
	AppID string
	// TimeStamp is the Unix timestamp string for the pay signature.
	TimeStamp string
	// NonceStr is a random string for the pay signature.
	NonceStr string
	// Package is the prepay_id package string (prepay_id=...).
	Package string
	// SignType is the signature algorithm (RSA).
	SignType string
	// PaySign is the RSA signature over the invoke parameters.
	PaySign string
}

// NativeResult holds the QR code URL.
type NativeResult struct {
	// CodeURL is the QR code URL for Native pay.
	CodeURL string
}

// H5Result holds the redirect URL.
type H5Result struct {
	// H5URL is the redirect URL for H5 pay.
	H5URL string
}

// APPResult holds prepay_id for the mobile app SDK.
type APPResult struct {
	// PrepayID is the prepay_id for the mobile SDK invoke.
	PrepayID string
}

// CombinePrepayResult holds combine prepay outcomes.
type CombinePrepayResult struct {
	// PrepayID is the combine prepay_id when applicable.
	PrepayID string
	// CodeURL is the Native QR URL when applicable.
	CodeURL string
	// H5URL is the H5 redirect URL when applicable.
	H5URL string
}

// PrepayResponse is filled according to Channel.
type PrepayResponse struct {
	// Channel echoes the prepay channel used.
	Channel Channel
	// PrepayID is the WeChat prepay_id when returned directly.
	PrepayID string

	// JSAPI holds front-end invoke params for JSAPI / mini program pay.
	JSAPI *JSAPIInvokeParams
	// Native holds the QR code URL for Native pay.
	Native *NativeResult
	// H5 holds the redirect URL for H5 pay.
	H5 *H5Result
	// APP holds prepay_id for the mobile SDK.
	APP *APPResult
	// Combine holds combine-specific prepay outcomes.
	Combine *CombinePrepayResult
}

// QueryOrderRequest queries by merchant order number.
type QueryOrderRequest struct {
	// OutTradeNo is the merchant-side order number.
	OutTradeNo string
	// Kind selects normal or combine order query API.
	Kind OrderKind
}

// CloseOrderRequest closes an unpaid order.
type CloseOrderRequest struct {
	// OutTradeNo is the merchant-side order number to close.
	OutTradeNo string
	// Kind selects normal or combine order close API.
	Kind OrderKind
}

// Transaction mirrors the paid order resource exposed to callers.
type Transaction struct {
	// AppID is the WeChat app ID.
	AppID string
	// MchID is the merchant ID.
	MchID string
	// OutTradeNo is the merchant-side order number.
	OutTradeNo string
	// TransactionID is the WeChat payment transaction ID.
	TransactionID string
	// TradeType is the trade type (JSAPI, NATIVE, etc.).
	TradeType string
	// TradeState is the order state (SUCCESS, NOTPAY, etc.).
	TradeState string
	// TradeStateDesc is a human-readable state description.
	TradeStateDesc string
	// Attach is merchant metadata from the order.
	Attach string
	// SuccessTime is RFC3339 payment success time.
	SuccessTime string
	// BankType is the payer's bank type when available.
	BankType string
	// Amount is the paid amount block.
	Amount *Amount
	// PayerOpenID is the payer's OpenID when available.
	PayerOpenID string
}

// RefundAmount is the refund amount block for Create refund.
type RefundAmount struct {
	// Total is the original order total in the smallest currency unit.
	Total int64
	// Refund is the refund amount in the smallest currency unit.
	Refund int64
	// Currency is the ISO 4217 code; defaults to CNY when empty.
	Currency string
}

// RefundRequest creates a domestic refund.
type RefundRequest struct {
	// OutTradeNo is the merchant-side order number being refunded.
	OutTradeNo string
	// OutRefundNo is the merchant-side unique refund number.
	OutRefundNo string
	// Reason is an optional refund reason shown to the payer.
	Reason string
	// NotifyURL receives refund result notifications.
	NotifyURL string
	// Amount is the refund amount block.
	Amount RefundAmount
}

// QueryRefundRequest queries refund by out_refund_no.
type QueryRefundRequest struct {
	// OutRefundNo is the merchant-side refund number.
	OutRefundNo string
}

// RefundResult is a domestic refund resource.
type RefundResult struct {
	// RefundID is the WeChat refund ID.
	RefundID string
	// OutRefundNo is the merchant-side refund number.
	OutRefundNo string
	// TransactionID is the original payment transaction ID.
	TransactionID string
	// OutTradeNo is the merchant-side order number.
	OutTradeNo string
	// Channel is the refund channel (ORIGINAL, BALANCE, etc.).
	Channel string
	// UserReceivedAccount is where the refund was credited.
	UserReceivedAccount string
	// SuccessTime is RFC3339 refund success time when completed.
	SuccessTime string
	// CreateTime is RFC3339 refund creation time.
	CreateTime string
	// Status is the refund status (SUCCESS, PROCESSING, etc.).
	Status string
	// Amount is the refund amount block.
	Amount *RefundAmount
}

// NotifySummary is the outer notification envelope after verification.
type NotifySummary struct {
	// ID is the notification ID from WeChat Pay.
	ID string
	// CreateTime is when the notification was created.
	CreateTime time.Time
	// EventType is the event type (e.g. TRANSACTION.SUCCESS).
	EventType string
	// ResourceType is the encrypted resource type.
	ResourceType string
	// Summary is a short human-readable summary.
	Summary string
}
