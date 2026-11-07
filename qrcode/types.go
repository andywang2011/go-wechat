package qrcode

// ImageResult is a successful wxacode / QR image response (binary body).
type ImageResult struct {
	// ContentType is the HTTP Content-Type of the image bytes.
	ContentType string
	// Data is the raw image payload.
	Data []byte
}

// LineColor for wxacode line color when auto_color is false.
type LineColor struct {
	// R is the red component (divided by 255 in the API).
	R int `json:"r"`
	// G is the green component.
	G int `json:"g"`
	// B is the blue component.
	B int `json:"b"`
}

// GetWxaCodeReq requests a limited wxacode for a specific path.
type GetWxaCodeReq struct {
	// Path is the mini program page path (max length per WeChat docs).
	Path string
	// Width is the QR image width in pixels.
	Width int
	// AutoColor uses default line colors when true.
	AutoColor bool
	// LineColor sets custom RGB when AutoColor is false.
	LineColor *LineColor
	// IsHyaline renders a transparent background when true.
	IsHyaline bool
}

// GetWxaCodeUnlimitReq requests unlimited wxacode (scene-based).
type GetWxaCodeUnlimitReq struct {
	// Scene is the scene parameter passed when scanning (max 32 chars).
	Scene string
	// Page is the target page; omit to use the homepage.
	Page string
	// CheckPath validates Page exists when true.
	CheckPath bool
	// EnvVersion is release, trial, or develop.
	EnvVersion string
	// Width is the QR image width in pixels.
	Width int
	// AutoColor uses default line colors when true.
	AutoColor bool
	// LineColor sets custom RGB when AutoColor is false.
	LineColor *LineColor
	// IsHyaline renders a transparent background when true.
	IsHyaline bool
}

// CreateWxaQRCodeReq requests a legacy mini program QR code.
type CreateWxaQRCodeReq struct {
	// Path is the mini program page path.
	Path string
	// Width is the QR image width in pixels (default 430).
	Width int
}

// JumpWxa is the target mini program page for scheme / link generation.
type JumpWxa struct {
	// Path is the page path without query.
	Path string `json:"path,omitempty"`
	// Query is the page query string without leading ?.
	Query string `json:"query,omitempty"`
	// EnvVersion is release, trial, or develop.
	EnvVersion string `json:"env_version,omitempty"`
}

// GenerateSchemeReq generates an encrypted URL scheme.
type GenerateSchemeReq struct {
	// JumpWxa is the landing page when the scheme is opened.
	JumpWxa *JumpWxa
	// IsExpire enables expiry when true.
	IsExpire bool
	// ExpireType is 0 (timestamp) or 1 (interval days).
	ExpireType int
	// ExpireTime is Unix timestamp expiry when ExpireType is 0.
	ExpireTime int64
	// ExpireInterval is days until expiry when ExpireType is 1.
	ExpireInterval int
}

// GenerateSchemeResp is the generatescheme response.
type GenerateSchemeResp struct {
	RespErr
	// OpenLink is the generated URL scheme.
	OpenLink string `json:"openlink,omitempty"`
}

// QuerySchemeReq queries scheme info.
type QuerySchemeReq struct {
	// Scheme is the URL scheme to query.
	Scheme string
}

// QuerySchemeResp is the queryscheme response.
type QuerySchemeResp struct {
	RespErr
	// SchemeInfo holds decoded scheme metadata.
	SchemeInfo *SchemeInfo `json:"scheme_info,omitempty"`
}

// SchemeInfo holds scheme query details.
type SchemeInfo struct {
	// Scheme is the URL scheme string.
	Scheme string `json:"scheme,omitempty"`
	// OpenLink is the equivalent open link.
	OpenLink string `json:"openlink,omitempty"`
	// ExpireTime is Unix expiry when applicable.
	ExpireTime int64 `json:"expire_time,omitempty"`
}

// GenerateNFCSchemeReq generates NFC scheme.
type GenerateNFCSchemeReq struct {
	// JumpWxa is the landing page when the NFC tag is tapped.
	JumpWxa *JumpWxa
	// SN is the device serial number.
	SN string
	// ModelID is the device model ID registered with WeChat.
	ModelID string
}

// GenerateNFCSchemeResp is the generatenfcscheme response.
type GenerateNFCSchemeResp struct {
	RespErr
	// OpenLink is the generated NFC scheme link.
	OpenLink string `json:"openlink,omitempty"`
}

// GenerateURLLinkReq generates encrypted URL Link.
type GenerateURLLinkReq struct {
	// Path is the mini program page path.
	Path string
	// Query is the page query string.
	Query string
	// EnvVersion is release, trial, or develop.
	EnvVersion string
	// ExpireType is 0 (timestamp) or 1 (interval days).
	ExpireType int
	// ExpireTime is Unix timestamp expiry when ExpireType is 0.
	ExpireTime int64
	// ExpireInterval is days until expiry when ExpireType is 1.
	ExpireInterval int
}

// GenerateURLLinkResp is the generate_urllink response.
type GenerateURLLinkResp struct {
	RespErr
	// URLLink is the generated URL Link.
	URLLink string `json:"url_link,omitempty"`
}

// QueryURLLinkReq queries URL Link config.
type QueryURLLinkReq struct {
	// URLLink is the URL Link to query.
	URLLink string
}

// QueryURLLinkResp is the query_urllink response.
type QueryURLLinkResp struct {
	RespErr
	// URLLinkInfo holds decoded URL Link metadata.
	URLLinkInfo *URLLinkInfo `json:"url_link_info,omitempty"`
}

// URLLinkInfo holds URL Link query details.
type URLLinkInfo struct {
	// URLLink is the URL Link string.
	URLLink string `json:"url_link,omitempty"`
	// ExpireTime is Unix expiry when applicable.
	ExpireTime int64 `json:"expire_time,omitempty"`
}

// GenShortLinkReq generates Short Link.
type GenShortLinkReq struct {
	// PageURL is the full mini program page URL.
	PageURL string
	// PageTitle is the link title shown when sharing.
	PageTitle string
	// IsPermanent creates a permanent link when true.
	IsPermanent bool
}

// GenShortLinkResp is the genwxashortlink response.
type GenShortLinkResp struct {
	RespErr
	// Link is the generated Short Link.
	Link string `json:"link,omitempty"`
}
