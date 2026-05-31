package application

import (
	"context"
	"time"

	"github.com/diagnosis/go-toolkit/logger"
)

func (a *Application) StartEventCleanup() {
	go func() {
		for {
			// Calculate time until next midnight
			now := time.Now()
			next := now.Add(24 * time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			timer := time.NewTimer(next.Sub(now))
			<-timer.C

			ctx := context.Background()
			err := a.eventStore.DeleteOldEvents(ctx)
			if err != nil {
				logger.Error(ctx, "failed to cleanup old events", "err", err)
			} else {
				logger.Info(ctx, "old events cleaned up")
			}
		}
	}()
}
