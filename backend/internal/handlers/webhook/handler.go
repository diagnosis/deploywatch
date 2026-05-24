package handlers

import (
	"github.com/diagnosis/deploywatchv2/internal/config"
	"github.com/diagnosis/deploywatchv2/internal/events"
	eventstore "github.com/diagnosis/deploywatchv2/internal/store/event"
	"github.com/diagnosis/deploywatchv2/internal/store/installation"
	userstore "github.com/diagnosis/deploywatchv2/internal/store/user"
	"github.com/diagnosis/deploywatchv2/internal/store/watchedrepo"
)

type WebhookHandler struct {
	cfg               *config.Config
	hub               *events.Hub
	eventStore        eventstore.EventStore
	watchedRepo       watchedrepo.WatchedRepoStore
	installationStore installation.InstallationStore
	userStore         userstore.UserStore
}

func NewWebhookHandler(
	cfg *config.Config,
	hub *events.Hub,
	eventStore eventstore.EventStore,
	watchedRepos watchedrepo.WatchedRepoStore,
	installationStore installation.InstallationStore,
	userStore userstore.UserStore,
) *WebhookHandler {
	return &WebhookHandler{
		cfg,
		hub,
		eventStore,
		watchedRepos,
		installationStore,
		userStore,
	}
}
