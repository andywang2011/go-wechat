package coupon

import (
	"errors"
	"fmt"

	libcommon "github.com/andywang2011/go-wechat/common"
)

// DefaultHost is the WeChat API host.
const DefaultHost = libcommon.DefaultAPIHost

// Config holds credentials for one application (keyed by appkey in the app layer).
type Config struct {
	// Host overrides DefaultHost when set.
	Host string
	// AppID is the WeChat application ID.
	AppID string
	// Secret is the WeChat application secret.
	Secret string
}

// ConfigProvider resolves config by appkey.
type ConfigProvider interface {
	GetConfig(appkey string) (Config, error)
}

// MapConfigProvider is a static map for tests or simple setups.
type MapConfigProvider map[string]Config

// GetConfig implements ConfigProvider.
func (m MapConfigProvider) GetConfig(appkey string) (Config, error) {
	if appkey == "" {
		return Config{}, errors.New("coupon: appkey is required")
	}
	cfg, ok := m[appkey]
	if !ok {
		return Config{}, fmt.Errorf("coupon: config not found for appkey %s", appkey)
	}
	if cfg.Host == "" {
		cfg.Host = DefaultHost
	}
	return cfg, nil
}
