package utils

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func Clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val >= max {
		return max - 1
	}
	return val
}
