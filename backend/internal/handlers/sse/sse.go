package handlers

import (
	"net/http"

	"github.com/diagnosis/deploywatchv2/internal/events"
	"github.com/diagnosis/go-toolkit/errors"
	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/diagnosis/go-toolkit/responder"
	"github.com/google/uuid"
)

func (h *SSEHandler) HandleSSE(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	correlationID, _ := logger.GetCorrelationID(ctx)

	idStr, ok := middleware.GetUserID(ctx)
	if !ok {
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		responder.Error(w, errors.Unauthorized("unauthorized", "unauthorized"), correlationID)
		return
	}
	writer, err := middleware.NewSSEWriter(w)
	if err != nil {
		responder.Error(w, err, correlationID)
		return
	}
	c := &events.Client{
		UserID: id,
		Send:   make(chan events.Event),
	}
	h.hub.Register(c)
	defer h.hub.Unregister(c)
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-c.Send:
			if !ok {
				return
			}
			if err := writer.Send(event.Type, event.Data); err != nil {
				return
			}

		}
	}
}
