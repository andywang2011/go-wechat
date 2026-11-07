package pay

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// AuthMode selects how outbound requests and notify verification authenticate.
type AuthMode int

const (
	// AuthModeAutoCert uses platform certificates auto-fetched by the SDK.
	AuthModeAutoCert AuthMode = iota
	// AuthModePublicKey uses the WeChat Pay merchant public key (recommended for new merchants).
	AuthModePublicKey
)

// Config holds merchant credentials injected by the application (no viper, no file paths).
type Config struct {
	// MchID is the WeChat Pay merchant ID.
	MchID string
	// MchCertSerialNo is the merchant API certificate serial number.
	MchCertSerialNo string
	// MchPrivateKeyPEM is the merchant API private key PEM.
	MchPrivateKeyPEM string
	// APIv3Key is the 32-byte API v3 key for decrypting callbacks.
	APIv3Key string
	// DefaultAppID is the default WeChat app ID for prepay when AppID is omitted.
	DefaultAppID string

	// AuthMode selects platform certificate auto-fetch or merchant public key auth.
	AuthMode AuthMode
	// WechatPayPublicKeyID is required when AuthMode is AuthModePublicKey.
	WechatPayPublicKeyID string
	// WechatPayPublicKeyPEM is required when AuthMode is AuthModePublicKey.
	WechatPayPublicKeyPEM string

	// Timeout is the HTTP client timeout for Pay API calls.
	Timeout time.Duration
}

// Validate checks required fields for the selected auth mode.
func (c Config) Validate() error {
	if strings.TrimSpace(c.MchID) == "" {
		return errors.New("pay: mch_id is required")
	}
	if strings.TrimSpace(c.MchCertSerialNo) == "" {
		return errors.New("pay: mch_cert_serial_no is required")
	}
	if strings.TrimSpace(c.MchPrivateKeyPEM) == "" {
		return errors.New("pay: mch_private_key_pem is required")
	}
	if _, err := utils.LoadPrivateKey(c.MchPrivateKeyPEM); err != nil {
		return fmt.Errorf("pay: invalid mch_private_key_pem: %w", err)
	}
	if strings.TrimSpace(c.APIv3Key) == "" {
		return errors.New("pay: api_v3_key is required")
	}
	if len(c.APIv3Key) != 32 {
		return errors.New("pay: api_v3_key must be 32 bytes")
	}
	switch c.AuthMode {
	case AuthModeAutoCert:
		// no extra fields
	case AuthModePublicKey:
		if strings.TrimSpace(c.WechatPayPublicKeyID) == "" {
			return errors.New("pay: wechat_pay_public_key_id is required for public key auth")
		}
		if strings.TrimSpace(c.WechatPayPublicKeyPEM) == "" {
			return errors.New("pay: wechat_pay_public_key_pem is required for public key auth")
		}
		if _, err := utils.LoadPublicKey(c.WechatPayPublicKeyPEM); err != nil {
			return fmt.Errorf("pay: invalid wechat_pay_public_key_pem: %w", err)
		}
	default:
		return fmt.Errorf("pay: unsupported auth mode %d", c.AuthMode)
	}
	return nil
}
