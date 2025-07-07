package tui

import (
	"fmt"
	"time"
)

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val >= max {
		return max - 1
	}
	return val
}
