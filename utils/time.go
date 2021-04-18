package utils

import (
	"time"
)

func IsExpired(startTime time.Time, period int) bool {
	now := time.Now().UTC().Truncate(time.Second)
	if now.Sub(startTime).Seconds() > float64(period) {
		return true
	}
	return false
}
