package pay

// Channel identifies which WeChat Pay v3 prepay endpoint to call.
type Channel string

const (
	// ChannelJSAPI JSAPI / Mini Program pay — POST /v3/pay/transactions/jsapi
	ChannelJSAPI Channel = "jsapi"
	// ChannelNative Native QR scan — POST /v3/pay/transactions/native
	ChannelNative Channel = "native"
	// ChannelH5 H5 pay — POST /v3/pay/transactions/h5
	ChannelH5 Channel = "h5"
	// ChannelAPP APP pay — POST /v3/pay/transactions/app
	ChannelAPP Channel = "app"

	// ChannelCombineJSAPI combine JSAPI — POST /v3/combine-transactions/jsapi
	ChannelCombineJSAPI Channel = "combine_jsapi"
	// ChannelCombineNative combine Native — POST /v3/combine-transactions/native
	ChannelCombineNative Channel = "combine_native"
	// ChannelCombineH5 combine H5 — POST /v3/combine-transactions/h5
	ChannelCombineH5 Channel = "combine_h5"
	// ChannelCombineAPP combine APP — POST /v3/combine-transactions/app
	ChannelCombineAPP Channel = "combine_app"
)

// APIPath returns the v3 path (without host) for prepay on this channel.
func (c Channel) APIPath() string {
	switch c {
	case ChannelJSAPI:
		return "/v3/pay/transactions/jsapi"
	case ChannelNative:
		return "/v3/pay/transactions/native"
	case ChannelH5:
		return "/v3/pay/transactions/h5"
	case ChannelAPP:
		return "/v3/pay/transactions/app"
	case ChannelCombineJSAPI:
		return "/v3/combine-transactions/jsapi"
	case ChannelCombineNative:
		return "/v3/combine-transactions/native"
	case ChannelCombineH5:
		return "/v3/combine-transactions/h5"
	case ChannelCombineAPP:
		return "/v3/combine-transactions/app"
	default:
		return ""
	}
}

// DocURL returns the official product documentation URL for the channel.
func (c Channel) DocURL() string {
	switch c {
	case ChannelJSAPI:
		return "https://pay.weixin.qq.com/docs/merchant/products/jsapi-payment/development.html"
	case ChannelNative:
		return "https://pay.weixin.qq.com/docs/merchant/products/native-payment/development.html"
	case ChannelH5:
		return "https://pay.weixin.qq.com/docs/merchant/products/h5-payment/development.html"
	case ChannelAPP:
		return "https://pay.weixin.qq.com/docs/merchant/products/in-app-payment/development.html"
	case ChannelCombineJSAPI, ChannelCombineNative, ChannelCombineH5, ChannelCombineAPP:
		return "https://pay.weixin.qq.com/docs/merchant/products/combine-payment/development.html"
	default:
		return ""
	}
}

// IsCombine reports whether the channel uses combine-transactions APIs.
func (c Channel) IsCombine() bool {
	switch c {
	case ChannelCombineJSAPI, ChannelCombineNative, ChannelCombineH5, ChannelCombineAPP:
		return true
	default:
		return false
	}
}

// IsValid reports whether c is a known channel constant.
func (c Channel) IsValid() bool {
	return c.APIPath() != ""
}
