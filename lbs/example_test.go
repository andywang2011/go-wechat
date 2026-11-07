package lbs_test

import (
	"fmt"
	"time"

	"github.com/andywang2011/go-wechat/lbs"
)

func ExampleNewClient() {
	_, err := lbs.NewClient(lbs.Config{AppKey: "your-tencent-lbs-key", Timeout: 10 * time.Second})
	if err != nil {
		panic(err)
	}
}

func Example_geocoderEndpoints() {
	fmt.Println(lbs.GeocoderV1URL)
	fmt.Println(lbs.ReverseGeocoderDocURL)
	// Output:
	// https://apis.map.qq.com/ws/geocoder/v1/
	// https://lbs.qq.com/service/webService/webServiceGuide/address/Gcoder
}
