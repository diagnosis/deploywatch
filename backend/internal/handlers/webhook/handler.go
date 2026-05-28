package handlers

import (
	"github.com/diagnosis/deploywatchv2/internal/config"
	"github.com/diagnosis/deploywatchv2/internal/events"
	"github.com/diagnosis/deploywatchv2/internal/fcm"
	devicetokenstore "github.com/diagnosis/deploywatchv2/internal/store/devicetoken"
	eventstore "github.com/diagnosis/deploywatchv2/internal/store/event"
	"github.com/diagnosis/deploywatchv2/internal/store/installation"
	userstore "github.com/diagnosis/deploywatchv2/internal/store/user"
	"github.com/diagnosis/deploywatchv2/internal/store/watchedrepo"
)

type WebhookHandler struct {
	cfg               *config.Config
	hub               *events.Hub
	wsHub             *events.WSHub
	eventStore        eventstore.EventStore
	watchedRepo       watchedrepo.WatchedRepoStore
	installationStore installation.InstallationStore
	userStore         userstore.UserStore
	fcm               *fcm.Client
	deviceTokenStore  devicetokenstore.DeviceTokenStore
}

func NewWebhookHandler(
	cfg *config.Config,
	hub *events.Hub,
	wsHub *events.WSHub,
	eventStore eventstore.EventStore,
	watchedRepos watchedrepo.WatchedRepoStore,
	installationStore installation.InstallationStore,
	userStore userstore.UserStore,
	fcm *fcm.Client,
	deviceTokenStore devicetokenstore.DeviceTokenStore,
) *WebhookHandler {
	return &WebhookHandler{
		cfg,
		hub,
		wsHub,
		eventStore,
		watchedRepos,
		installationStore,
		userStore,
		fcm,
		deviceTokenStore,
	}
}
