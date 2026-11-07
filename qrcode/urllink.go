package qrcode

import "context"

func (c *Client) GenerateURLLink(ctx context.Context, appkey, accessToken string, req GenerateURLLinkReq) (*GenerateURLLinkResp, error) {
	body := map[string]interface{}{}
	if req.Path != "" {
		body["path"] = req.Path
	}
	if req.Query != "" {
		body["query"] = req.Query
	}
	if req.EnvVersion != "" {
		body["env_version"] = req.EnvVersion
	}
	if req.ExpireType > 0 || req.ExpireTime > 0 || req.ExpireInterval > 0 {
		body["expire_type"] = req.ExpireType
		if req.ExpireTime > 0 {
			body["expire_time"] = req.ExpireTime
		}
		if req.ExpireInterval > 0 {
			body["expire_interval"] = req.ExpireInterval
		}
	}
	var resp GenerateURLLinkResp
	err := c.postJSON(ctx, appkey, accessToken, PathGenerateURLLink, body, &resp)
	if err != nil {
		return nil, err
	}
	if err := resp.RespErr.error(); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) QueryURLLink(ctx context.Context, appkey, accessToken string, req QueryURLLinkReq) (*QueryURLLinkResp, error) {
	var resp QueryURLLinkResp
	err := c.postJSON(ctx, appkey, accessToken, PathQueryURLLink, map[string]string{
		"url_link": req.URLLink,
	}, &resp)
	if err != nil {
		return nil, err
	}
	if err := resp.RespErr.error(); err != nil {
		return nil, err
	}
	return &resp, nil
}
