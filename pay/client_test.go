package pay

import "testing"

func TestNewClient_ValidateError(t *testing.T) {
	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewClientWithGateway(t *testing.T) {
	cfg := validTestConfig(t)
	api := NewClientWithGateway(cfg, &mockGateway{}, nil)
	if api == nil {
		t.Fatal("expected non-nil API")
	}
}
