package handlers

import (
	"github.com/diagnosis/deploywatchv2/internal/apns"
	"github.com/diagnosis/deploywatchv2/internal/config"
	"github.com/diagnosis/deploywatchv2/internal/events"
	devicetokenstore "github.com/diagnosis/deploywatchv2/internal/store/devicetoken"
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
	apnsClient        *apns.Client
	deviceTokenStore  devicetokenstore.DeviceTokenStore
}

func NewWebhookHandler(
	cfg *config.Config,
	hub *events.Hub,
	eventStore eventstore.EventStore,
	watchedRepos watchedrepo.WatchedRepoStore,
	installationStore installation.InstallationStore,
	userStore userstore.UserStore,
	apnsClient *apns.Client,
	deviceTokenStore devicetokenstore.DeviceTokenStore,
) *WebhookHandler {
	return &WebhookHandler{
		cfg,
		hub,
		eventStore,
		watchedRepos,
		installationStore,
		userStore,
		apnsClient,
		deviceTokenStore,
	}
}
