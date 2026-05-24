package handlers

import (
	"net/http"
	"strconv"

	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *WatchedRepoHandler) HandleRemoveRepo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	userIdStr, ok := middleware.GetUserID(ctx)
	if !ok {
		logger.Error(ctx, "unauthorized")
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		logger.Error(ctx, "unauthorized", "err", err)
		responder.Error(w, errors.Unauthorized("unauthorized", "failed to parse id"), correlationID)
		return
	}

	idStr := chi.URLParam(r, "id")
	repoID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error(ctx, "bad repo id", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	if err := h.watchRepoStore.Remove(ctx, userID, repoID); err != nil {
		logger.Error(ctx, "failed to remove repo", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	w.WriteHeader(http.StatusOK)

}
