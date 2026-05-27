package handlers

import (
	"net/http"

	"github.com/diagnosis/deploywatchv2/internal/events"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	userIDStr, ok := middleware.GetUserID(ctx)
	if !ok {
		responder.Error(w, errors.Unauthorized("unauthoerized", "unauthorized"), correlationID)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}

	websocket.Handler(func(conn *websocket.Conn) {
		client := &events.WSClient{
			UserID: userID,
			Send:   make(chan events.Event, 32),
		}

		h.wsHub.Register(client)

		defer h.wsHub.Unregister(client)

		logger.Info(ctx, "ws client connected", "user_id", userID)

		for event := range client.Send {
			msg := event.Type + "\n" + event.Data
			if err := websocket.Message.Send(conn, msg); err != nil {
				logger.Error(ctx, "ws send failed", "err", err)
				break
			}
		}
	}).ServeHTTP(w, r)
}
