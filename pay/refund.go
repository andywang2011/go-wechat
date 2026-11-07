package pay

import (
	"fmt"
	"strings"
)

func (r RefundRequest) validate() error {
	if strings.TrimSpace(r.OutRefundNo) == "" {
		return fmt.Errorf("pay: out_refund_no is required")
	}
	if strings.TrimSpace(r.OutTradeNo) == "" {
		return fmt.Errorf("pay: out_trade_no is required")
	}
	if r.Amount.Refund <= 0 || r.Amount.Total <= 0 {
		return fmt.Errorf("pay: refund amount total and refund must be positive")
	}
	return nil
}

func (r QueryRefundRequest) validate() error {
	if strings.TrimSpace(r.OutRefundNo) == "" {
		return fmt.Errorf("pay: out_refund_no is required")
	}
	return nil
}
