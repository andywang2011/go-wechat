package qrcode_test

import (
	"context"
	"fmt"

	"github.com/andywang2011/go-wechat/qrcode"
)

func ExampleClient_GenerateScheme() {
	client := qrcode.NewClient(qrcode.MapConfigProvider{"app1": {}})
	accessToken := "..." // from token center
	resp, err := client.GenerateScheme(context.Background(), "app1", accessToken, qrcode.GenerateSchemeReq{
		JumpWxa: &qrcode.JumpWxa{Path: "pages/index", Query: "id=1"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.OpenLink)
}

func ExampleDocURL() {
	fmt.Println(qrcode.PathGetWxaCodeUnlimit)
	fmt.Println(qrcode.DocURL(qrcode.PathGetWxaCodeUnlimit))
	// Output:
	// /wxa/getwxacodeunlimit
	// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/qrcode-link/qr-code/getUnlimitedQRCode.html
}
