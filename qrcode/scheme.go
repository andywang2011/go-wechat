package qrcode

import "context"

func (c *Client) QueryScheme(ctx context.Context, appkey, accessToken string, req QuerySchemeReq) (*QuerySchemeResp, error) {
	var resp QuerySchemeResp
	err := c.postJSON(ctx, appkey, accessToken, PathQueryScheme, map[string]string{
		"scheme": req.Scheme,
	}, &resp)
	if err != nil {
		return nil, err
	}
	if err := resp.RespErr.error(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GenerateScheme(ctx context.Context, appkey, accessToken string, req GenerateSchemeReq) (*GenerateSchemeResp, error) {
	body := map[string]interface{}{}
	if req.JumpWxa != nil {
		body["jump_wxa"] = req.JumpWxa
	}
	if req.IsExpire {
		body["is_expire"] = true
		body["expire_type"] = req.ExpireType
		if req.ExpireTime > 0 {
			body["expire_time"] = req.ExpireTime
		}
		if req.ExpireInterval > 0 {
			body["expire_interval"] = req.ExpireInterval
		}
	}
	var resp GenerateSchemeResp
	err := c.postJSON(ctx, appkey, accessToken, PathGenerateScheme, body, &resp)
	if err != nil {
		return nil, err
	}
	if err := resp.RespErr.error(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) GenerateNFCScheme(ctx context.Context, appkey, accessToken string, req GenerateNFCSchemeReq) (*GenerateNFCSchemeResp, error) {
	body := map[string]interface{}{
		"sn":       req.SN,
		"model_id": req.ModelID,
	}
	if req.JumpWxa != nil {
		body["jump_wxa"] = req.JumpWxa
	}
	var resp GenerateNFCSchemeResp
	err := c.postJSON(ctx, appkey, accessToken, PathGenerateNFCScheme, body, &resp)
	if err != nil {
		return nil, err
	}
	if err := resp.RespErr.error(); err != nil {
		return nil, err
	}
	return &resp, nil
}
