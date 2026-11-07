package pay

import (
	"fmt"
	"strings"
)

func (r QueryOrderRequest) validate() error {
	if strings.TrimSpace(r.OutTradeNo) == "" {
		return fmt.Errorf("pay: out_trade_no is required")
	}
	return nil
}

func (r CloseOrderRequest) validate() error {
	if strings.TrimSpace(r.OutTradeNo) == "" {
		return fmt.Errorf("pay: out_trade_no is required")
	}
	return nil
}
