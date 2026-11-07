package pay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func generateTestKeyPair(t *testing.T) (privatePEM, publicPEM string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	privDER, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	privatePEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER}))
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	return privatePEM, publicPEM
}

func validTestConfig(t *testing.T) Config {
	t.Helper()
	priv, pub := generateTestKeyPair(t)
	return Config{
		MchID:                 "1900000001",
		MchCertSerialNo:       "serial",
		MchPrivateKeyPEM:      priv,
		APIv3Key:              "12345678901234567890123456789012",
		DefaultAppID:          "wx_appid",
		AuthMode:              AuthModePublicKey,
		WechatPayPublicKeyID:  "PUBKEYID",
		WechatPayPublicKeyPEM: pub,
	}
}
