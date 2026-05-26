package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
)

type registerTokenReq struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

func (h *DeviceTokenHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	userIDStr, ok := middleware.GetUserID(ctx)
	if !ok {
		responder.Error(w, errors.Unauthorized("unauthorized", "no user id"), correlationID)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responder.Error(w, errors.BadRequest("invalid user id", "invalid user id"), correlationID)
		return
	}

	var req registerTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		responder.Error(w, errors.BadRequest("invalid body", "token required"), correlationID)
		return
	}

	if req.Platform == "" {
		req.Platform = "ios"
	}

	if err := h.store.Upsert(ctx, userID, req.Token, req.Platform); err != nil {
		logger.Error(ctx, "failed to upsert device token", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	responder.JSON(w, http.StatusOK, map[string]string{"status": "ok"}, correlationID)
}
