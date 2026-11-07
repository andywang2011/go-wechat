package user_test

import (
	"context"
	"fmt"

	"github.com/andywang2011/go-wechat/user"
)

func ExampleClient_GetUserPhoneNumber() {
	client := user.NewClient(user.MapConfigProvider{
		"app1": {AppID: "wx_appid", Secret: "app_secret"},
	})
	accessToken := "..." // from application token center
	resp, err := client.GetUserPhoneNumber(context.Background(), "app1", accessToken, "phone_code")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.PhoneInfo.PurePhoneNumber)
}

func ExampleGetUserPhoneNumberPath() {
	fmt.Println(user.GetUserPhoneNumberPath)
	// Output: /wxa/business/getuserphonenumber
}
