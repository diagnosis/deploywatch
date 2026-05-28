package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/diagnosis/deploywatchv2/internal/events"
	eventstore "github.com/diagnosis/deploywatchv2/internal/store/event"

	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/diagnosis/go-toolkit/secure"
)

type githubPayload struct {
	Action     string `json:"action"`
	Repository struct {
		ID int64 `json:"id"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}

	sigHeader := r.Header.Get("X-Hub-Signature-256")
	if sigHeader == "" {
		responder.Error(w, errors.BadRequest("missing signature", "missing signature"), correlationID)
		return
	}
	sigHeader = strings.TrimPrefix(sigHeader, "sha256=")
	sigByte, err := hex.DecodeString(sigHeader)
	if err != nil {
		logger.Error(ctx, "invalid sig", "err", err)
		responder.Error(w, errors.BadRequest("invalid signature format", "invalid signature format"), correlationID)
		return
	}
	if ok := secure.VerifyHMACSHA256([]byte(h.cfg.GitHub.GitHubWebhookSecret), b, sigByte); !ok {
		responder.Error(w, errors.BadRequest("wrong sig", "wrong sig"), correlationID)
		return
	}

	var payload githubPayload
	if err := json.Unmarshal(b, &payload); err != nil {
		logger.Error(ctx, "failed to parse payload", "err", err)
	}
	action := payload.Action
	var actionPtr *string
	if action != "" {
		actionPtr = &action
	}

	eventType := r.Header.Get("X-GitHub-Event")
	logger.Info(ctx, "webhook received", "event", eventType)

	_, err = h.eventStore.Create(ctx, &eventstore.Event{
		RepoID:     payload.Repository.ID,
		EventType:  eventType,
		Action:     actionPtr,
		ActorLogin: payload.Sender.Login,
		Payload:    b,
	})
	if err != nil {
		logger.Error(ctx, "failed to save event", "err", err)
	}
	if eventType == "installation" {
		if err := h.HandleInstallation(ctx, b); err != nil {
			logger.Error(ctx, "failed to handle installation", "err", err)
		}
	}

	userIDs, err := h.watchedRepo.GetUsersByRepoID(ctx, payload.Repository.ID)
	if err != nil {
		logger.Error(ctx, "failed to get watching users", "err", err)

	}
	logger.Info(ctx, "broadcasting to users", "count", len(userIDs), "repo_id", payload.Repository.ID)

	for _, userID := range userIDs {
		h.hub.BroadcastToUser(userID, events.Event{
			Type: eventType,
			Data: string(b),
		})
		h.wsHub.BroadcastToUser(userID, events.Event{
			Type: eventType,
			Data: string(b),
		})
	}

	if h.fcm != nil {
		for _, userID := range userIDs {
			tokens, err := h.deviceTokenStore.GetByUserID(ctx, userID)
			if err != nil {
				logger.Error(ctx, "failed to get device tokens", "err", err)
				continue
			}
			for _, t := range tokens {
				go func(token string) {
					fcmCtx := context.Background()
					title := eventType
					body := fmt.Sprintf("%s by %s", payload.Repository.ID, payload.Sender.Login)
					if err := h.fcm.Send(fcmCtx, token, title, body); err != nil {
						logger.Error(ctx, " failed to send push", "err", err)
					}

				}(t.Token)
			}
		}

	}

	responder.JSON(w, http.StatusOK, nil, correlationID)

}
