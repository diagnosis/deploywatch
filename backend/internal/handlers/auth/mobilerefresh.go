package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	refreshtokenstore "github.com/diagnosis/deploywatchv2/internal/store/refreshtoken"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/diagnosis/go-toolkit/secure"
)

func (h *AuthHandler) HandleMobileRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		responder.Error(w, errors.BadRequest("invalid body", "refresh_token required"), correlationID)
		return
	}

	hash := secure.HashRefreshToken(body.RefreshToken)
	rt, err := h.refreshTokenStore.GetByTokenHash(ctx, hash)
	if err != nil || time.Now().After(rt.ExpiresAt) {
		responder.Error(w, errors.Unauthorized("invalid token", "refresh token invalid or expired"), correlationID)
		return
	}

	token, err := h.jwt.SignAccess(rt.UserID.String())
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}

	err = h.refreshTokenStore.DeleteByUserIDAndPlatform(ctx, rt.UserID, rt.Platform)
	if err != nil {
		logger.Warn(ctx, "failed to delete refresh tokens", "err", err)
	}
	raw, _ := secure.GenerateRefreshToken()
	newHash := secure.HashRefreshToken(raw)
	err = h.refreshTokenStore.Create(ctx, &refreshtokenstore.RefreshToken{
		UserID:    rt.UserID,
		TokenHash: newHash,
		Platform:  rt.Platform,
		ExpiresAt: time.Now().Add(h.cfg.JWT.RefreshTokenExpiry),
		CreatedAt: time.Now(),
	})
	if err != nil {
		logger.Warn(ctx, "failed to create refresh token", "err", err)
	}

	responder.JSON(w, http.StatusOK, map[string]string{
		"access_token":  token,
		"refresh_token": raw,
	}, correlationID)
}
