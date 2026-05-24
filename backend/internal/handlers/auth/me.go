package handlers

import (
	"net/http"

	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
)

func (h *AuthHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)
	//getting userID
	idStr, ok := middleware.GetUserID(ctx)
	if !ok {
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error(ctx, "failed to parse id", "err", err)
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}
	//fetch user
	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		logger.Error(ctx, "failed to fetch user", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	responder.JSON(w, http.StatusOK, map[string]any{
		"user": user,
	}, correlationID)
}
