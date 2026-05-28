package handlers

import (
	"net/http"

	"github.com/diagnosis/deploywatchv2/internal/events"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins for now
	},
}

func (h *WSHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	userIDStr, ok := middleware.GetUserID(ctx)
	if !ok {
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(ctx, "ws upgrade failed", "err", err)
		return
	}
	defer conn.Close()

	client := &events.WSClient{
		UserID: userID,
		Send:   make(chan events.Event, 32),
	}

	h.wsHub.Register(client)
	defer h.wsHub.Unregister(client)

	logger.Info(ctx, "ws client connected", "user_id", userID)

	for event := range client.Send {
		msg := event.Type + "\n" + event.Data
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			logger.Error(ctx, "ws send failed", "err", err)
			break
		}
	}
}
