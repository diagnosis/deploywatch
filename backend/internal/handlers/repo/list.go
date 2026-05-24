package handlers

import (
	"net/http"

	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
)

func (h *WatchedRepoHandler) HandleListRepo(w http.ResponseWriter, r *http.Request) {
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

	repos, err := h.watchRepoStore.GetByUserID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "failed to fetch repos", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	count := len(repos)

	responder.JSON(w, http.StatusOK, map[string]any{
		"repos": repos,
		"count": count,
	}, correlationID)
}
