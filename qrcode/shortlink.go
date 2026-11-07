package qrcode

import "context"

func (c *Client) GenShortLink(ctx context.Context, appkey, accessToken string, req GenShortLinkReq) (*GenShortLinkResp, error) {
	body := map[string]interface{}{
		"page_url": req.PageURL,
	}
	if req.PageTitle != "" {
		body["page_title"] = req.PageTitle
	}
	if req.IsPermanent {
		body["is_permanent"] = true
	}
	var resp GenShortLinkResp
	err := c.postJSON(ctx, appkey, accessToken, PathGenShortLink, body, &resp)
	if err != nil {
		return nil, err
	}
	if err := resp.RespErr.error(); err != nil {
		return nil, err
	}
	return &resp, nil
}
