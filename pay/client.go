package pay

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

// Client implements API using wechatpay-go under the hood.
type Client struct {
	cfg    Config
	gw     paymentGateway
	notify *notify.Handler
}

// NewClient builds a production client from Config.
func NewClient(cfg Config) (API, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	privateKey, err := utils.LoadPrivateKey(cfg.MchPrivateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("pay: load private key: %w", err)
	}
	opts := []core.ClientOption{
		option.WithMerchantCredential(cfg.MchID, cfg.MchCertSerialNo, privateKey),
	}
	var verifier auth.Verifier
	switch cfg.AuthMode {
	case AuthModeAutoCert:
		opts = append(opts, option.WithWechatPayAutoAuthCipher(
			cfg.MchID, cfg.MchCertSerialNo, privateKey, cfg.APIv3Key,
		))
		certVisitor := downloader.MgrInstance().GetCertificateVisitor(cfg.MchID)
		verifier = verifiers.NewSHA256WithRSAVerifier(certVisitor)
	case AuthModePublicKey:
		publicKey, err := utils.LoadPublicKey(cfg.WechatPayPublicKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("pay: load public key: %w", err)
		}
		opts = append(opts, option.WithWechatPayPublicKeyAuthCipher(
			cfg.MchID, cfg.MchCertSerialNo, privateKey, cfg.WechatPayPublicKeyID, publicKey,
		))
		verifier = verifiers.NewSHA256WithRSAPubkeyVerifier(cfg.WechatPayPublicKeyID, *publicKey)
	default:
		return nil, fmt.Errorf("pay: unsupported auth mode %d", cfg.AuthMode)
	}
	if cfg.Timeout > 0 {
		opts = append(opts, option.WithHTTPClient(&http.Client{Timeout: cfg.Timeout}))
	}
	coreClient, err := core.NewClient(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("pay: create core client: %w", err)
	}
	notifyHandler, err := notify.NewRSANotifyHandler(cfg.APIv3Key, verifier)
	if err != nil {
		return nil, fmt.Errorf("pay: create notify handler: %w", err)
	}
	return &Client{
		cfg:    cfg,
		gw:     newSDKGateway(coreClient),
		notify: notifyHandler,
	}, nil
}

// NewClientWithGateway injects a gateway for tests.
func NewClientWithGateway(cfg Config, gw paymentGateway, notifyHandler *notify.Handler) API {
	return &Client{cfg: cfg, gw: gw, notify: notifyHandler}
}

func (c *Client) Prepay(ctx context.Context, req PrepayRequest) (*PrepayResponse, error) {
	if err := req.Validate(c.cfg); err != nil {
		return nil, err
	}
	return c.gw.prepay(ctx, req, c.cfg)
}

func (c *Client) QueryOrder(ctx context.Context, req QueryOrderRequest) (*Transaction, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return c.gw.queryOrder(ctx, req, c.cfg)
}

func (c *Client) CloseOrder(ctx context.Context, req CloseOrderRequest) error {
	if err := req.validate(); err != nil {
		return err
	}
	return c.gw.closeOrder(ctx, req, c.cfg)
}

func (c *Client) Refund(ctx context.Context, req RefundRequest) (*RefundResult, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return c.gw.refund(ctx, req, c.cfg)
}

func (c *Client) QueryRefund(ctx context.Context, req QueryRefundRequest) (*RefundResult, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return c.gw.queryRefund(ctx, req, c.cfg)
}
