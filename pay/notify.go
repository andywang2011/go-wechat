package pay

import (
	"context"
	"fmt"
	"net/http"
)

// ParseNotify verifies the signature, decrypts the resource, and unmarshals into dst.
func (c *Client) ParseNotify(ctx context.Context, r *http.Request, dst any) (*NotifySummary, error) {
	if c.notify == nil {
		return nil, fmt.Errorf("pay: notify handler not configured")
	}
	if dst == nil {
		return nil, fmt.Errorf("pay: dst is required")
	}
	nreq, err := c.notify.ParseNotifyRequest(ctx, r, dst)
	if err != nil {
		return nil, err
	}
	summary := &NotifySummary{
		ID:           nreq.ID,
		EventType:    nreq.EventType,
		ResourceType: nreq.ResourceType,
		Summary:      nreq.Summary,
	}
	if nreq.CreateTime != nil {
		summary.CreateTime = *nreq.CreateTime
	}
	return summary, nil
}
