package pay

import "testing"

func TestChannel_APIPath(t *testing.T) {
	cases := []struct {
		ch   Channel
		want string
	}{
		{ChannelJSAPI, "/v3/pay/transactions/jsapi"},
		{ChannelNative, "/v3/pay/transactions/native"},
		{ChannelCombineJSAPI, "/v3/combine-transactions/jsapi"},
	}
	for _, tc := range cases {
		if got := tc.ch.APIPath(); got != tc.want {
			t.Errorf("%s APIPath = %q, want %q", tc.ch, got, tc.want)
		}
	}
}

func TestChannel_IsCombine(t *testing.T) {
	if !ChannelCombineH5.IsCombine() {
		t.Error("combine h5 should be combine")
	}
	if ChannelJSAPI.IsCombine() {
		t.Error("jsapi should not be combine")
	}
}

func TestChannel_DocURL(t *testing.T) {
	if ChannelJSAPI.DocURL() == "" {
		t.Error("jsapi doc url should be set")
	}
}

func TestChannel_IsValid(t *testing.T) {
	if Channel("unknown").IsValid() {
		t.Error("unknown channel should be invalid")
	}
}
