// /backend/internal/application/app.go
package application

import (
	"context"
	"fmt"

	"github.com/diagnosis/deploywatchv2/internal/apns"
	"github.com/diagnosis/deploywatchv2/internal/config"
	"github.com/diagnosis/deploywatchv2/internal/database"
	"github.com/diagnosis/deploywatchv2/internal/events"
	"github.com/diagnosis/deploywatchv2/internal/github"
	authhandler "github.com/diagnosis/deploywatchv2/internal/handlers/auth"
	devicetokenhandler "github.com/diagnosis/deploywatchv2/internal/handlers/devicetoken"
	eventhandler "github.com/diagnosis/deploywatchv2/internal/handlers/event"
	githubclienthandlers "github.com/diagnosis/deploywatchv2/internal/handlers/github"
	healthhandler "github.com/diagnosis/deploywatchv2/internal/handlers/health"
	watchedrepohandlers "github.com/diagnosis/deploywatchv2/internal/handlers/repo"
	ssehandlers "github.com/diagnosis/deploywatchv2/internal/handlers/sse"
	webhookhandler "github.com/diagnosis/deploywatchv2/internal/handlers/webhook"
	wshandlers "github.com/diagnosis/deploywatchv2/internal/handlers/ws"
	devicetokenstore "github.com/diagnosis/deploywatchv2/internal/store/devicetoken"
	eventstore "github.com/diagnosis/deploywatchv2/internal/store/event"

	"github.com/diagnosis/deploywatchv2/internal/store/installation"
	refreshtokenstore "github.com/diagnosis/deploywatchv2/internal/store/refreshtoken"
	userstore "github.com/diagnosis/deploywatchv2/internal/store/user"
	"github.com/diagnosis/deploywatchv2/internal/store/watchedrepo"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/secure"
)

type Application struct {
	// jwtSigner
	jwt          *secure.JWTSigner
	Hub          *events.Hub
	WsHub        *events.WSHub
	gitHubClient *github.GitHubClient
	APNSClient   *apns.Client
	// stores
	userStore         userstore.UserStore
	refreshTokenStore refreshtokenstore.RefreshTokenStore
	eventStore        eventstore.EventStore
	watchedRepos      watchedrepo.WatchedRepoStore
	installationStore installation.InstallationStore
	deviceTokenStore  devicetokenstore.DeviceTokenStore
	//handlers
	healthHandler       *healthhandler.HealthHandler
	authHandler         *authhandler.AuthHandler
	webhookHandler      *webhookhandler.WebhookHandler
	sseHandler          *ssehandlers.SSEHandler
	wsHandler           *wshandlers.WSHandler
	watchedRepoHandler  *watchedrepohandlers.WatchedRepoHandler
	eventHandler        *eventhandler.EventHandler
	gitHubClientHandler *githubclienthandlers.GitHubHandler
	deviceTokenHandler  *devicetokenhandler.DeviceTokenHandler
}

func NewApplication() (*Application, error) {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger.Init()
	pool, err := database.OpenPool(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	hub := events.NewHub()
	wsHub := events.NewWSHub()
	gitHubClient := github.NewGitHubClient(cfg)
	APNSClient, err := apns.New(cfg.APNS)
	if err != nil {
		logger.Error(ctx, "failed to fetch apns client", "err", err)
	}

	jwt, err := secure.NewJWTSigner(secure.JWTConfig{
		AccessSecret:       cfg.JWT.AccessSecret,
		RefreshSecret:      cfg.JWT.RefreshSecret,
		AccessTokenExpiry:  cfg.JWT.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
		Issuer:             cfg.JWT.Issuer,
		Audience:           cfg.JWT.Audience,
		Leeway:             0, //default 30s
	})
	if err != nil {
		logger.Fatal(ctx, "failed to set admin jwt signer", "err", err)
	}
	// stores
	userStore := userstore.NewPGUserStore(pool)
	refreshTokenStore := refreshtokenstore.NewPGRefreshTokenStore(pool)
	eventStore := eventstore.NewPGEventStore(pool)
	watchedRepoStore := watchedrepo.NewPGWatchedRepoStore(pool)
	installationStore := installation.NewPGInstallationStore(pool)
	deviceTokenStore := devicetokenstore.NewPGDeviceTokenStore(pool)
	// handlers
	healthHandler := healthhandler.NewHealthHandler()
	authHandler := authhandler.NewAuthHandler(cfg, jwt, userStore, refreshTokenStore)
	webhookHandler := webhookhandler.NewWebhookHandler(
		cfg,
		hub,
		wsHub,
		eventStore,
		watchedRepoStore,
		installationStore,
		userStore,
		APNSClient,
		deviceTokenStore,
	)
	sseHandler := ssehandlers.NewSSEHandler(hub)
	wsHandler := wshandlers.NewWSHandler(wsHub)
	watchedRepoHandler := watchedrepohandlers.NewWatchRepoHandler(watchedRepoStore, installationStore)
	eventHandler := eventhandler.NewEventHandler(eventStore, watchedRepoStore)
	gitHubClientHandler := githubclienthandlers.NewGitHubHandler(gitHubClient, installationStore)
	deviceTokenHandler := devicetokenhandler.NewDeviceTokenHandler(deviceTokenStore)
	return &Application{
		jwt:               jwt,
		Hub:               hub,
		WsHub:             wsHub,
		gitHubClient:      gitHubClient,
		APNSClient:        APNSClient,
		userStore:         userStore,
		refreshTokenStore: refreshTokenStore,
		eventStore:        eventStore,
		installationStore: installationStore,
		watchedRepos:      watchedRepoStore,
		deviceTokenStore:  deviceTokenStore,

		healthHandler:       healthHandler,
		authHandler:         authHandler,
		webhookHandler:      webhookHandler,
		sseHandler:          sseHandler,
		wsHandler:           wsHandler,
		watchedRepoHandler:  watchedRepoHandler,
		eventHandler:        eventHandler,
		gitHubClientHandler: gitHubClientHandler,
		deviceTokenHandler:  deviceTokenHandler,
	}, nil
}
