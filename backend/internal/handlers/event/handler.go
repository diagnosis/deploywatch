package handlers

import (
	eventstore "github.com/diagnosis/deploywatchv2/internal/store/event"
	"github.com/diagnosis/deploywatchv2/internal/store/watchedrepo"
)

type EventHandler struct {
	eventStore       eventstore.EventStore
	watchedRepoStore watchedrepo.WatchedRepoStore
}

func NewEventHandler(eventStore eventstore.EventStore, watchedRepoStore watchedrepo.WatchedRepoStore) *EventHandler {
	return &EventHandler{
		eventStore:       eventStore,
		watchedRepoStore: watchedRepoStore,
	}
}
