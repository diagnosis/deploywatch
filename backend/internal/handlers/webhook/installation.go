package handlers

import (
	"context"
	"encoding/json"

	"github.com/diagnosis/deploywatchv2/internal/store/installation"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
)

type installationPayload struct {
	Action       string `json:"action"`
	Installation struct {
		ID      int64 `json:"id"`
		Account struct {
			Login string `json:"login"`
			Type  string `json:"type"`
		} `json:"account"`
	} `json:"installation"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`
}

func (h *WebhookHandler) HandleInstallation(ctx context.Context, payload []byte) error {
	var ip installationPayload
	err := json.Unmarshal(payload, &ip)
	if err != nil {
		logger.Error(ctx, "failed to parse payload", "err", err)
		return err
	}

	switch ip.Action {
	case "deleted":
		err := h.installationStore.Delete(ctx, ip.Installation.ID)
		if err != nil {
			return err
		}
	case "created":
		user, err := h.userStore.GetUserByLogin(ctx, ip.Sender.Login)
		if err != nil {
			return err
		}
		_, err = h.installationStore.Upsert(ctx, &installation.Installation{
			InstallationID: ip.Installation.ID,
			UserID:         user.ID,
			AccountLogin:   user.Login,
			AccountType:    ip.Installation.Account.Type,
		})
		if err != nil {
			return err
		}
	case "suspend":
		err := h.installationStore.Suspend(ctx, ip.Installation.ID)
		if err != nil {
			return err
		}
	case "unsuspend":
		err := h.installationStore.Unsuspend(ctx, ip.Installation.ID)
		if err != nil {
			return err
		}
	default:
		return errors.BadRequest("action not allowed", "action not allowed")
	}
	return nil

}
