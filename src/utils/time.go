package utils

import "fmt"

func FormatDuration(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	var result string
	if days > 0 {
		result += fmt.Sprintf("%dd, ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh, ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm, ", minutes)
	}
	result += fmt.Sprintf("%ds", remainingSeconds)

	return result
}
