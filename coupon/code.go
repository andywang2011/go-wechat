package coupon

import "context"

const (
	// GetCodePath POST /card/code/get
	// Doc: https://developers.weixin.qq.com/doc/offiaccount/Cards_and_Offer/Managing_Coupons_Vouchers_and_Cards.html
	GetCodePath = "/card/code/get"
)

// GetCodeReq is the card code lookup request.
type GetCodeReq struct {
	Code         string `json:"code"`
	CardID       string `json:"card_id,omitempty"`
	CheckConsume bool   `json:"check_consume,omitempty"`
}

// GetCodeResp is the card code lookup response.
type GetCodeResp struct {
	RespErr
	CardID string `json:"card_id,omitempty"`
	OpenID string `json:"openid,omitempty"`
}

// GetCode obtains card code metadata (requires access_token).
func (c *Client) GetCode(ctx context.Context, appkey, accessToken string, req *GetCodeReq) (*GetCodeResp, error) {
	var resp GetCodeResp
	err := c.post(ctx, appkey, accessToken, GetCodePath, req, &resp)
	return &resp, err
}
