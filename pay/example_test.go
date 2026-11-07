package pay_test

import (
	"fmt"

	"github.com/andywang2011/go-wechat/pay"
)

// Production: read certificate files in the application and inject PEM strings (library does not read paths).
//
//	client, err := pay.NewClient(pay.Config{
//	    MchID:            os.Getenv("WECHAT_MCH_ID"),
//	    MchCertSerialNo:  os.Getenv("WECHAT_MCH_SERIAL"),
//	    MchPrivateKeyPEM: string(privateKeyPEM),
//	    APIv3Key:         os.Getenv("WECHAT_API_V3_KEY"),
//	    DefaultAppID:     os.Getenv("WECHAT_APP_ID"),
//	    AuthMode:         pay.AuthModePublicKey,
//	    WechatPayPublicKeyID:  os.Getenv("WECHAT_PAY_PUBKEY_ID"),
//	    WechatPayPublicKeyPEM: string(wechatPayPublicKeyPEM),
//	})
func ExampleNewClient() {
	cfg := pay.Config{
		MchID:            "1900000001",
		MchCertSerialNo:  "serial",
		MchPrivateKeyPEM: "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
		APIv3Key:         "12345678901234567890123456789012",
		DefaultAppID:     "wx_appid",
	}
	_, err := pay.NewClient(cfg)
	if err != nil {
		fmt.Println("configure valid PEM in production")
		return
	}
}

// Channel constants include v3 paths and official documentation links.
func ExampleChannel_metadata() {
	fmt.Println(pay.ChannelJSAPI.APIPath())
	fmt.Println(pay.ChannelJSAPI.DocURL())
	// Output:
	// /v3/pay/transactions/jsapi
	// https://pay.weixin.qq.com/docs/merchant/products/jsapi-payment/development.html
}
