package apns

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/diagnosis/deploywatchv2/internal/config"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

type Client struct {
	client   *apns2.Client
	bundleID string
}

func New(cfg config.APNSConfig) (*Client, error) {
	keyData, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read APNs key: %w", err)
	}

	authKey, err := authKeyFromBytes(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse APNs key: %w", err)
	}

	t := &token.Token{
		AuthKey: authKey,
		KeyID:   cfg.KeyID,
		TeamID:  cfg.TeamID,
	}

	client := apns2.NewTokenClient(t).Production()

	return &Client{
		client:   client,
		bundleID: cfg.BundleID,
	}, nil
}

func (c *Client) Send(deviceToken string, title string, body string) error {
	notification := &apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       c.bundleID,
		Payload: map[string]any{
			"aps": map[string]any{
				"alert": map[string]any{
					"title": title,
					"body":  body,
				},
				"sound": "default",
			},
		},
	}

	res, err := c.client.Push(notification)
	if err != nil {
		return fmt.Errorf("APNs push failed: %w", err)
	}
	if !res.Sent() {
		return fmt.Errorf("APNs push not sent: %s", res.Reason)
	}
	return nil
}

func authKeyFromBytes(bytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not ECDSA")
	}
	return ecKey, nil
}
