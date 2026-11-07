package common

// DefaultAPIHost is the default WeChat Open Platform API host.
const DefaultAPIHost = "https://api.weixin.qq.com"

// AppCredential holds appid/secret for a WeChat application.
type AppCredential struct {
	// Host overrides DefaultAPIHost when set.
	Host string
	// AppID is the WeChat application ID.
	AppID string
	// Secret is the WeChat application secret.
	Secret string
}

// HostOrDefault returns Host or DefaultAPIHost when empty.
func (c AppCredential) HostOrDefault() string {
	if c.Host != "" {
		return c.Host
	}
	return DefaultAPIHost
}
