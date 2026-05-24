package handlers

import (
	"net/http"

	"github.com/diagnosis/go-toolkit/logger"
	"github.com/diagnosis/go-toolkit/responder"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	correlationID, _ := logger.GetCorrelationID(r.Context())
	responder.JSON(w, http.StatusOK, map[string]string{
		"message": "server is ok",
	}, correlationID)
}
