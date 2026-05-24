package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/diagnosis/deploywatchv2/internal/store/watchedrepo"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
)

type WatchRepoReq struct {
	RepoID         int64  `json:"repo_id"`
	RepoFullName   string `json:"repo_full_name"`
	InstallationID int64  `json:"installation_id"`
}

func (h *WatchedRepoHandler) HandleAddRepo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	correlationID, _ := logger.GetCorrelationID(ctx)

	idStr, ok := middleware.GetUserID(ctx)
	if !ok {
		logger.Error(ctx, "unauthorized")
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error(ctx, "unauthorized", "err", err)
		responder.Error(w, errors.Unauthorized("unauthorized", "failed to parse id"), correlationID)
		return
	}

	var req WatchRepoReq
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&req)
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}

	repo, err := h.watchRepoStore.Add(ctx, &watchedrepo.WatchedRepo{
		UserID:         id,
		InstallationID: req.InstallationID,
		RepoID:         req.RepoID,
		RepoFullName:   req.RepoFullName,
		EventFilters:   []byte(`{}`),
	})
	if err != nil {
		logger.Error(ctx, "failed to add repo", "err", err)
		responder.Error(w, err, correlationID)
		return
	}

	responder.JSON(w, http.StatusCreated, map[string]any{
		"repo": repo,
	}, correlationID)
}
