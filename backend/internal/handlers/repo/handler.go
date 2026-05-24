package handlers

import (
	"github.com/diagnosis/deploywatchv2/internal/store/installation"
	"github.com/diagnosis/deploywatchv2/internal/store/watchedrepo"
)

type WatchedRepoHandler struct {
	watchRepoStore    watchedrepo.WatchedRepoStore
	installationStore installation.InstallationStore
}

func NewWatchRepoHandler(
	watchedRepoStore watchedrepo.WatchedRepoStore,
	installationStore installation.InstallationStore) *WatchedRepoHandler {
	return &WatchedRepoHandler{
		watchRepoStore:    watchedRepoStore,
		installationStore: installationStore,
	}
}
