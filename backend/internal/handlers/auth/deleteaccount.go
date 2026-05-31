package handlers

import (
	"net/http"

	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
)

func (h *AuthHandler) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	idStr, ok := middleware.GetUserID(ctx)
	if !ok {
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}

	err = h.userStore.DeleteUser(ctx, id)
	if err != nil {
		logger.Error(ctx, "failed to delete user", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	h.clearRefreshCookie(w)
	h.clearAccessCookie(w)
	responder.JSON(w, http.StatusOK, nil, correlationID)
}
