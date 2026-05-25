package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	refreshtokenstore "github.com/diagnosis/deploywatchv2/internal/store/refreshtoken"
	userstore "github.com/diagnosis/deploywatchv2/internal/store/user"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/diagnosis/go-toolkit/secure"
	"golang.org/x/oauth2"
)

func (h *AuthHandler) HandleCallBack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	stateToken := r.URL.Query().Get("state")
	if stateToken == "" {
		logger.Error(ctx, "bad state token")
		responder.Error(w, errors.BadRequest("bad token", "bad state token"), correlationID)
		return
	}

	stateTokenCookie, err := r.Cookie("state_token")
	if err != nil {
		logger.Error(ctx, "error getting state token cookie", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	if stateToken != stateTokenCookie.Value {
		logger.Error(ctx, "bad state token")
		responder.Error(w, errors.BadRequest("bad token", "bad state token"), correlationID)
		return
	}

	// code

	code := r.URL.Query().Get("code")
	accessToken, err := h.exchangeCodeForToken(ctx, code)
	if err != nil {
		logger.Error(ctx, "error exchanging code for token", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken}))
	githubUser, err := secure.FetchGitHubUserInfo(ctx, oauthClient)
	if err != nil {
		logger.Error(ctx, "failed to fetch gh user info", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	user, err := h.userStore.UpsertUser(ctx, &userstore.User{
		GitHubID:    githubUser.ID,
		Login:       githubUser.Login,
		Name:        &githubUser.Name,
		AvatarURL:   &githubUser.AvatarURL,
		Email:       &githubUser.Email,
		AccessToken: accessToken,
	})
	if err != nil {
		logger.Error(ctx, "failed to upsert user", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	//
	raw, err := secure.GenerateRefreshToken()
	if err != nil {
		logger.Error(ctx, "failed to generate refresh token", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	hash := secure.HashRefreshToken(raw)
	err = h.refreshTokenStore.Create(ctx, &refreshtokenstore.RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(h.cfg.JWT.RefreshTokenExpiry),
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.Error(ctx, "failed to save refresh token in db", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	token, err := h.jwt.SignAccess(user.ID.String())
	if err != nil {
		logger.Error(ctx, "failed to mint access token", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	// Mobile flow
	mobileCookie, err := r.Cookie("mobile_redirect")
	if err == nil {
		mobileRedirect := fmt.Sprintf("%s?access_token=%s&refresh_token=%s", mobileCookie.Value, token, raw)
		http.Redirect(w, r, mobileRedirect, http.StatusTemporaryRedirect)
		return
	}

	h.setRefreshCookie(w, raw)
	h.setAccessCookie(w, token)

	http.Redirect(w, r, h.cfg.App.FrontendURL, http.StatusTemporaryRedirect)

}

type GithubClientRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}
type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func (h *AuthHandler) exchangeCodeForToken(ctx context.Context, code string) (string, error) {
	payload := GithubClientRequest{
		ClientID:     h.cfg.GitHub.GitHubOauthClientID,
		ClientSecret: h.cfg.GitHub.GitHubOauthClientSecret,
		Code:         code,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(
		ctx, "POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github token exchange failed: status %d", resp.StatusCode)
	}
	var ghTokenResp GithubTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&ghTokenResp)
	if err != nil {
		return "", err
	}

	return ghTokenResp.AccessToken, nil
}

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, token string) {

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.jwt.RefreshExpiry().Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}
func (h *AuthHandler) setAccessCookie(w http.ResponseWriter, token string) {

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(h.jwt.AccessExpiry().Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}
func (h *AuthHandler) clearAccessCookie(w http.ResponseWriter) {

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   int(-1),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}
func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   int(-1),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}
