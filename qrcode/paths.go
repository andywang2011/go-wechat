package qrcode

// API path constants with official documentation URLs (qrcode-link section).
const (
	// PathGetWxaCode POST /wxa/getwxacode — get Mini Program code (limited quantity)
	PathGetWxaCode = "/wxa/getwxacode"
	// PathGetWxaCodeUnlimit POST /wxa/getwxacodeunlimit — get unlimited Mini Program code
	PathGetWxaCodeUnlimit = "/wxa/getwxacodeunlimit"
	// PathCreateWxaQRCode POST /cgi-bin/wxaapp/createwxaqrcode — get Mini Program QR code
	PathCreateWxaQRCode = "/cgi-bin/wxaapp/createwxaqrcode"

	// PathQueryScheme POST /wxa/queryscheme — query URL scheme
	PathQueryScheme = "/wxa/queryscheme"
	// PathGenerateScheme POST /wxa/generatescheme — generate encrypted URL scheme
	PathGenerateScheme = "/wxa/generatescheme"
	// PathGenerateNFCScheme POST /wxa/generatenfcscheme — generate NFC URL scheme
	PathGenerateNFCScheme = "/wxa/generatenfcscheme"

	// PathGenerateURLLink POST /wxa/generate_urllink — generate encrypted URL Link
	PathGenerateURLLink = "/wxa/generate_urllink"
	// PathQueryURLLink POST /wxa/query_urllink — query encrypted URL Link
	PathQueryURLLink = "/wxa/query_urllink"

	// PathGenShortLink POST /wxa/genwxashortlink — generate Short Link
	PathGenShortLink = "/wxa/genwxashortlink"
)

const docBase = "https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link"

// DocURL returns the official documentation URL for an API path.
func DocURL(path string) string {
	switch path {
	case PathGetWxaCode:
		return docBase + "/qr-code/getQRCode.html"
	case PathGetWxaCodeUnlimit:
		return docBase + "/qr-code/getUnlimitedQRCode.html"
	case PathCreateWxaQRCode:
		return docBase + "/qr-code/getQRCode.html"
	case PathQueryScheme:
		return docBase + "/url-scheme/queryScheme.html"
	case PathGenerateScheme:
		return docBase + "/url-scheme/generateScheme.html"
	case PathGenerateNFCScheme:
		return docBase + "/url-scheme/generateNFCScheme.html"
	case PathGenerateURLLink:
		return docBase + "/url-link/generateUrlLink.html"
	case PathQueryURLLink:
		return docBase + "/url-link/queryUrlLink.html"
	case PathGenShortLink:
		return docBase + "/short-link/generateShortLink.html"
	default:
		return docBase + "/qr-code/getQRCode.html"
	}
}
