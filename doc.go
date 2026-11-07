// Package wechat provides SDKs for Tencent / WeChat product APIs.
//
// Design constraints (for use in application frameworks):
//   - Credentials and endpoints come from caller-injected Config; no viper or log.Fatal in the library.
//   - NewClient-style constructors return errors; callers own lifecycle and access_token caching via common.
//   - IO methods accept context.Context; no business logic in the SDK layer.
//   - ConfigDefault: N/A — each subpackage defines its own Config struct.
//
// Subpackages:
//
//   - common: shared access_token fetch and AccessTokenProvider contract
//   - lbs: Tencent Location Service (geocoding)
//   - login: WeChat Mini Program login (code2Session, checkSession, resetSession)
//   - user: WeChat Mini Program user info (phone number, etc.)
//   - qrcode: Mini Program QR code, URL Scheme, URL Link, Short Link
//   - subscribe: subscribe message and shared mini program config provider
//   - pay: WeChat Pay (reserved)
//   - coupon: WeChat coupon / card (reserved)
package wechat
