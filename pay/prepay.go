package pay

import (
	"fmt"
	"strings"
	"time"
)

const (
	outTradeNoMinLen = 6
	outTradeNoMaxLen = 32
)

// Validate checks PrepayRequest for the selected channel.
func (r PrepayRequest) Validate(cfg Config) error {
	if !r.Channel.IsValid() {
		return fmt.Errorf("pay: invalid channel %q", r.Channel)
	}
	if err := r.Base.validate(); err != nil {
		return err
	}
	if r.Channel.IsCombine() {
		return r.validateCombine(cfg)
	}
	switch r.Channel {
	case ChannelJSAPI:
		if r.Payer == nil || strings.TrimSpace(r.Payer.OpenID) == "" {
			return fmt.Errorf("pay: payer openid is required for %s", r.Channel)
		}
	case ChannelH5:
		if r.Scene == nil || strings.TrimSpace(r.Scene.Type) == "" {
			return fmt.Errorf("pay: h5 scene type is required for %s", r.Channel)
		}
	case ChannelNative, ChannelAPP:
		// base fields only
	}
	if resolveAppID(r.Base, cfg) == "" {
		return fmt.Errorf("pay: app_id is required (set base.app_id or config.default_app_id)")
	}
	return nil
}

func (b PrepayBase) validate() error {
	if strings.TrimSpace(b.Description) == "" {
		return fmt.Errorf("pay: description is required")
	}
	n := len(b.OutTradeNo)
	if n < outTradeNoMinLen || n > outTradeNoMaxLen {
		return fmt.Errorf("pay: out_trade_no length must be %d-%d", outTradeNoMinLen, outTradeNoMaxLen)
	}
	if strings.TrimSpace(b.NotifyURL) == "" {
		return fmt.Errorf("pay: notify_url is required")
	}
	if !strings.HasPrefix(strings.ToLower(b.NotifyURL), "https://") {
		return fmt.Errorf("pay: notify_url must use HTTPS")
	}
	if b.Amount.Total <= 0 {
		return fmt.Errorf("pay: amount.total must be positive")
	}
	if b.TimeExpire != "" {
		if _, err := time.Parse(time.RFC3339, b.TimeExpire); err != nil {
			return fmt.Errorf("pay: time_expire must be RFC3339: %w", err)
		}
	}
	return nil
}

func (r PrepayRequest) validateCombine(cfg Config) error {
	if r.Combine == nil {
		return fmt.Errorf("pay: combine info is required for %s", r.Channel)
	}
	if strings.TrimSpace(r.Combine.CombineOutTradeNo) == "" {
		return fmt.Errorf("pay: combine_out_trade_no is required")
	}
	if len(r.Combine.SubOrders) == 0 {
		return fmt.Errorf("pay: at least one sub_order is required")
	}
	for i, sub := range r.Combine.SubOrders {
		if strings.TrimSpace(sub.MchID) == "" {
			return fmt.Errorf("pay: sub_orders[%d].mch_id is required", i)
		}
		if strings.TrimSpace(sub.Description) == "" {
			return fmt.Errorf("pay: sub_orders[%d].description is required", i)
		}
		if strings.TrimSpace(sub.OutTradeNo) == "" {
			return fmt.Errorf("pay: sub_orders[%d].out_trade_no is required", i)
		}
		if sub.Amount.Total <= 0 {
			return fmt.Errorf("pay: sub_orders[%d].amount.total must be positive", i)
		}
	}
	switch r.Channel {
	case ChannelCombineJSAPI:
		if r.Payer == nil || strings.TrimSpace(r.Payer.OpenID) == "" {
			return fmt.Errorf("pay: payer openid is required for %s", r.Channel)
		}
	case ChannelCombineH5:
		if r.Scene == nil || strings.TrimSpace(r.Scene.Type) == "" {
			return fmt.Errorf("pay: h5 scene type is required for %s", r.Channel)
		}
	}
	if resolveAppID(r.Base, cfg) == "" {
		return fmt.Errorf("pay: app_id is required")
	}
	return nil
}
