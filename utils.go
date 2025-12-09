package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pterm/pterm"
	"strings"
)

func GetCountryOrigin(ip string) (IPStat, error) {

	url := fmt.Sprintf("http://ip-api.com/json/%v", ip)

	res, err := http.Get(url)

	if err != nil {
		return IPStat{}, err
	}

	defer res.Body.Close()

	stat := IPStat{}

	json.NewDecoder(res.Body).Decode(&stat)

	return stat, nil

}

func GetTotalDuration(sessions []AttemptSession) int {
	totalSeconds := 0

	for _, session := range sessions {
		var dur float64
		if !session.End.IsZero() {
			dur = session.End.Sub(session.Start).Seconds()
		} else {
			dur = time.Since(session.Start).Seconds()
		}
		if dur < 0 {
			dur = 0
		}
		totalSeconds += int(dur)
	}

	return totalSeconds
}

func GetWarmth(sessions []AttemptSession) (uint8, uint8, uint8) {
	if len(sessions) == 0 {
		return 0, 0, 0
	}

	totalSeconds := GetTotalDuration(sessions)

	if totalSeconds < 60 {
		return 255, 165, 0 // Orange
	} else if totalSeconds < 300 {
		return 255, 255, 0 // Yellow
	} else {
		return 255, 0, 0 // Red
	}

}

func RGBify(r, g, b uint8, row []string) []string {
	coloured := []string{}

	for _, value := range row {
		coloured = append(coloured, pterm.NewRGB(r, g, b).Sprint(value))
	}

	return coloured

}

func FormatDuration(totalSeconds int) string {
	if totalSeconds == 0 {
		return "0 seconds"
	}

	weeks := totalSeconds / (7 * 24 * 3600)
	totalSeconds %= (7 * 24 * 3600)
	days := totalSeconds / (24 * 3600)
	totalSeconds %= (24 * 3600)
	hours := totalSeconds / 3600
	totalSeconds %= 3600
	minutes := totalSeconds / 60
	totalSeconds %= 60
	seconds := totalSeconds

	var parts []string

	if weeks > 0 {
		parts = append(parts, fmt.Sprintf("%d weeks", weeks))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minutes", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%d seconds", seconds))
	}

	return strings.Join(parts, " ")
}
