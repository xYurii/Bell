package utils

import (
	"fmt"
	"strings"
)

func ConvertStringToSeconds(time string) int {
	time = strings.ToLower(time)
	parts := strings.Fields(time)
	totalSeconds := 0

	for _, part := range parts {
		value := 0
		unit := ""

		for i, char := range part {
			if char >= '0' && char <= '9' {
				value = value*10 + int(char-'0')
			} else {
				unit = part[i:]
				break
			}
		}

		switch unit {
		case "d", "dia", "dias":
			totalSeconds += value * 86400
		case "h", "hora", "horas":
			totalSeconds += value * 3600
		case "m", "minuto", "minutos":
			totalSeconds += value * 60
		case "s", "segundo", "segundos":
			totalSeconds += value
		}
	}

	return totalSeconds
}

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
