package fcm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Client struct {
	projectID   string
	tokenSource oauth2.TokenSource
}

type fcmMessage struct {
	Message message `json:"message"`
}

type message struct {
	Token        string       `json:"token"`
	Notification notification `json:"notification"`
	APNS         apnsConfig   `json:"apns"`
}
type apnsConfig struct {
	Payload apnsPayload `json:"payload"`
}
type apnsPayload struct {
	APS apsDict `json:"aps"`
}
type apsDict struct {
	Badge int `json:"badge"`
}

type notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func New(projectID, keyPath string) (*Client, error) {
	ctx := context.Background()

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read FCM key: %w", err)
	}

	creds, err := google.CredentialsFromJSONWithTypeAndParams(ctx, keyData,
		google.ServiceAccount,
		google.CredentialsParams{
			Scopes: []string{"https://www.googleapis.com/auth/firebase.messaging"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FCM credentials: %w", err)
	}

	return &Client{
		projectID:   projectID,
		tokenSource: creds.TokenSource,
	}, nil
}

func (c *Client) Send(ctx context.Context, deviceToken, title, body string) error {
	token, err := c.tokenSource.Token()
	if err != nil {
		return fmt.Errorf("failed to get FCM token: %w", err)
	}

	msg := fcmMessage{
		Message: message{
			Token: deviceToken,
			Notification: notification{
				Title: title,
				Body:  body,
			},
			APNS: apnsConfig{
				Payload: apnsPayload{
					APS: apsDict{Badge: 1},
				},
			},
		},
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal FCM message: %w", err)
	}

	url := fmt.Sprintf(
		"https://fcm.googleapis.com/v1/projects/%s/messages:send",
		c.projectID,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("failed to create FCM request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("FCM request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FCM returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
