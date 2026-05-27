package handlers

import "github.com/diagnosis/deploywatchv2/internal/events"

type WSHandler struct {
	wsHub *events.WSHub
}

func NewWSHandler(wsHub *events.WSHub) *WSHandler {
	return &WSHandler{wsHub}
}
