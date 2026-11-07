package subscribe

const (
	// SendSubscriptPath POST /cgi-bin/message/subscribe/send
	// Doc: https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/sendMessage.html
	SendSubscriptPath = "/cgi-bin/message/subscribe/send"
)

// SubscribeDataItem is one template field value in a subscribe message.
type SubscribeDataItem struct {
	// Value is the rendered text for the template field.
	Value string `json:"value"`
}

// SendSubscriptReq is a subscribe message request.
type SendSubscriptReq struct {
	// Touser is the recipient OpenID.
	Touser string `json:"touser"`
	// TemplateID is the subscribe message template ID.
	TemplateID string `json:"template_id"`
	// Page is the mini program page path opened after the user taps the message.
	Page string `json:"page,omitempty"`
	// Data maps template field names to rendered values.
	Data map[string]SubscribeDataItem `json:"data"`
	// MiniprogramState is developer, trial, or formal (default formal).
	MiniprogramState string `json:"miniprogram_state,omitempty"`
	// Lang is the language tag (default zh_CN).
	Lang string `json:"lang,omitempty"`
}

// SendSubscriptResp is the subscribe send response.
type SendSubscriptResp struct {
	RespErr
}
