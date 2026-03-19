package ui

import (
	"fmt"
	"time"

	"github.com/jasonsoprovich/shellclock/internal/model"
)

// isIdleWarning returns true when the active timer has been running (not
// paused) past the store's configured idle-warn threshold and the feature is
// enabled. Returns false when there is no timer or the timer is paused.
func isIdleWarning(store *model.Store) bool {
	if !store.IdleWarn.Enabled {
		return false
	}
	at := store.ActiveTimer
	if at == nil || at.Paused {
		return false
	}
	elapsed := at.AccumulatedSeconds + int64(time.Since(at.Start).Seconds())
	threshold := int64(store.IdleWarn.ThresholdMins) * 60
	return elapsed >= threshold
}

// idleWarnLabel returns a short warning string showing how many full hours the
// timer has been running, e.g. "⚠ 2h+".
func idleWarnLabel(store *model.Store) string {
	at := store.ActiveTimer
	if at == nil {
		return ""
	}
	elapsed := at.AccumulatedSeconds
	if !at.Paused {
		elapsed += int64(time.Since(at.Start).Seconds())
	}
	h := elapsed / 3600
	return fmt.Sprintf("⚠ %dh+", h)
}
