package login

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

const sigMethodHMACSHA256 = "hmac_sha256"

// SessionSignature returns login-state signature: hmac_sha256(session_key, "") as lowercase hex.
//
// Doc: https://developers.weixin.qq.com/miniprogram/dev/server/API/user-login/api_checksessionkey.html
func SessionSignature(sessionKey string) string {
	mac := hmac.New(sha256.New, []byte(sessionKey))
	_, _ = mac.Write(nil)
	return hex.EncodeToString(mac.Sum(nil))
}
