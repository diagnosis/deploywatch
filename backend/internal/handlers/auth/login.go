package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/diagnosis/go-toolkit/secure"
)

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	stateToken, err := secure.GenerateStateToken()
	if err != nil {
		logger.Error(ctx, "failed to generate token", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	h.setStateTokenCookie(w, stateToken)

	githubURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&state=%s&scope=read:user,user:email",
		h.cfg.GitHub.GitHubOauthClientID,
		stateToken,
	)
	http.Redirect(w, r, githubURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) setStateTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "state_token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(15 * time.Minute.Seconds()),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}
