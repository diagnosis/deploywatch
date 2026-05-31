package handlers

import (
	"net/http"
	"time"

	store "github.com/diagnosis/deploywatchv2/internal/store/refreshtoken"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/diagnosis/go-toolkit/secure"
)

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	// validate if refresh_token is valid
	c, err := r.Cookie("refresh_token")
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}
	hash := secure.HashRefreshToken(c.Value)
	refreshToken, err := h.refreshTokenStore.GetByTokenHash(ctx, hash)
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		responder.Error(w, errors.TokenError("expired token", "expired token"), correlationID)
		return
	}

	// mint new refresh and access token
	t, err := secure.GenerateRefreshToken()
	if err != nil {
		logger.Error(ctx, "failed to generate token", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	err = h.refreshTokenStore.DeleteByUserIDAndPlatform(ctx, refreshToken.UserID, "web")
	if err != nil {
		logger.Warn(ctx, "failed to delete old tokens", "err", err)
	}

	err = h.refreshTokenStore.Create(ctx, &store.RefreshToken{
		UserID:    refreshToken.UserID,
		TokenHash: secure.HashRefreshToken(t),
		Platform:  "web",
		ExpiresAt: time.Now().Add(h.cfg.JWT.RefreshTokenExpiry),
	})
	if err != nil {
		logger.Warn(ctx, "failed to save refresh token", "err", err)
	}
	at, err := h.jwt.SignAccess(refreshToken.UserID.String())
	if err != nil {
		logger.Error(ctx, "failed to generate token", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	h.setRefreshCookie(w, t)
	h.setAccessCookie(w, at)

	responder.JSON(w, http.StatusOK, nil, correlationID)

}
