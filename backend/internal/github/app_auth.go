package github

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (g *GitHubClient) GenerateAppJWT() (string, error) {
	block, _ := pem.Decode(g.cfg.GitHub.GitHubAppPrivateKey)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Add(-60 * time.Second).Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": g.cfg.GitHub.GitHubAppID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (g *GitHubClient) GetInstallationToken(ctx context.Context, installationID int64) (string, error) {
	// generate appJWT
	appjwt, err := g.GenerateAppJWT()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationID),
		nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+appjwt)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "deploywatch")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to get installation token: status %d", resp.StatusCode)
	}
	var tr tokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tr)
	if err != nil {
		return "", err
	}
	return tr.Token, nil

}

type tokenResponse struct {
	Token string `json:"token"`
}
