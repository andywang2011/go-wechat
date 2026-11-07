package qrcode

import "context"

func (c *Client) GetWxaCode(ctx context.Context, appkey, accessToken string, req GetWxaCodeReq) (*ImageResult, error) {
	body := map[string]interface{}{
		"path": req.Path,
	}
	if req.Width > 0 {
		body["width"] = req.Width
	}
	if req.AutoColor {
		body["auto_color"] = true
	}
	if req.LineColor != nil {
		body["line_color"] = req.LineColor
	}
	if req.IsHyaline {
		body["is_hyaline"] = true
	}
	return c.postImage(ctx, appkey, accessToken, PathGetWxaCode, body)
}

func (c *Client) GetWxaCodeUnlimit(ctx context.Context, appkey, accessToken string, req GetWxaCodeUnlimitReq) (*ImageResult, error) {
	body := map[string]interface{}{
		"scene": req.Scene,
	}
	if req.Page != "" {
		body["page"] = req.Page
	}
	if req.CheckPath {
		body["check_path"] = true
	}
	if req.EnvVersion != "" {
		body["env_version"] = req.EnvVersion
	}
	if req.Width > 0 {
		body["width"] = req.Width
	}
	if req.AutoColor {
		body["auto_color"] = true
	}
	if req.LineColor != nil {
		body["line_color"] = req.LineColor
	}
	if req.IsHyaline {
		body["is_hyaline"] = true
	}
	return c.postImage(ctx, appkey, accessToken, PathGetWxaCodeUnlimit, body)
}

func (c *Client) CreateWxaQRCode(ctx context.Context, appkey, accessToken string, req CreateWxaQRCodeReq) (*ImageResult, error) {
	body := map[string]interface{}{
		"path": req.Path,
	}
	if req.Width > 0 {
		body["width"] = req.Width
	}
	return c.postImage(ctx, appkey, accessToken, PathCreateWxaQRCode, body)
}
