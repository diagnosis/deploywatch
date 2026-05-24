package handlers

import (
	"net/http"
	"strconv"

	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
)

func (h *EventHandler) HandleListEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

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

	repos, err := h.watchedRepoStore.GetByUserID(ctx, id)
	if err != nil {
		logger.Error(ctx, "failed to get watched repos", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	var repoIDs []int64
	for _, repo := range repos {
		repoIDs = append(repoIDs, repo.RepoID)
	}

	repoIDStr := r.URL.Query().Get("repo_id")
	if repoIDStr != "" {
		repoID, err := strconv.ParseInt(repoIDStr, 10, 64)
		if err == nil {
			repoIDs = []int64{repoID}
		}
	}

	events, err := h.eventStore.ListByRepoIDs(ctx, repoIDs)
	if err != nil {
		logger.Error(ctx, "failed to get events", "err", err)
		responder.Error(w, err, correlationID)
		return
	}
	count := len(events)
	responder.JSON(w, http.StatusOK, map[string]any{
		"events": events,
		"count":  count,
	}, correlationID)
}
