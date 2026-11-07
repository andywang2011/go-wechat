package user

import "context"

const (
	// GetUserPhoneNumberPath POST /wxa/business/getuserphonenumber
	// Doc: https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-info/phone-number/getPhoneNumber.html
	GetUserPhoneNumberPath = "/wxa/business/getuserphonenumber"
)

// GetUserPhoneNumberResp is the getPhoneNumber response.
type GetUserPhoneNumberResp struct {
	RespErr
	PhoneInfo PhoneInfo `json:"phone_info,omitempty"`
}

// PhoneInfo holds the decrypted phone number fields.
type PhoneInfo struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

// GetUserPhoneNumber obtains the user's phone number from the dynamic code (requires access_token).
func (c *Client) GetUserPhoneNumber(ctx context.Context, appkey, accessToken, code string) (*GetUserPhoneNumberResp, error) {
	var resp GetUserPhoneNumberResp
	err := c.post(ctx, appkey, accessToken, GetUserPhoneNumberPath,
		map[string]string{"code": code}, &resp)
	return &resp, err
}
