package pay

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestClient_ParseNotify_NoHandler(t *testing.T) {
	cfg := validTestConfig(t)
	c := NewClientWithGateway(cfg, &mockGateway{}, nil).(*Client)
	_, err := c.ParseNotify(context.Background(), httptest.NewRequest("POST", "/", nil), &struct{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClient_ParseNotify_NilDst(t *testing.T) {
	cfg := validTestConfig(t)
	c := &Client{cfg: cfg, notify: nil, gw: &mockGateway{}}
	c.notify = nil
	_, err := c.ParseNotify(context.Background(), httptest.NewRequest("POST", "/", nil), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
