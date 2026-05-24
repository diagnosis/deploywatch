package handlers

import (
	"github.com/diagnosis/deploywatchv2/internal/events"
)

type SSEHandler struct {
	hub *events.Hub
}

func NewSSEHandler(hub *events.Hub) *SSEHandler {
	return &SSEHandler{hub}
}
