package pay

import (
	"testing"
)

func TestConfig_Validate_PublicKey(t *testing.T) {
	if err := validTestConfig(t).Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestConfig_Validate_MissingMchID(t *testing.T) {
	cfg := validTestConfig(t)
	cfg.MchID = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Validate_InvalidPrivatePEM(t *testing.T) {
	cfg := validTestConfig(t)
	cfg.MchPrivateKeyPEM = "not-a-key"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Validate_AutoCert(t *testing.T) {
	cfg := validTestConfig(t)
	cfg.AuthMode = AuthModeAutoCert
	cfg.WechatPayPublicKeyID = ""
	cfg.WechatPayPublicKeyPEM = ""
	if err := cfg.Validate(); err != nil {
		t.Fatalf("AutoCert Validate: %v", err)
	}
}

func TestConfig_Validate_PublicKeyMissingID(t *testing.T) {
	cfg := validTestConfig(t)
	cfg.WechatPayPublicKeyID = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Validate_APIv3KeyLength(t *testing.T) {
	cfg := validTestConfig(t)
	cfg.APIv3Key = "short"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}
