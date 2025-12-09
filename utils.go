package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pterm/pterm"
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

func GetWarmth(sessions []AttemptSession) (uint8, uint8, uint8) {
	if len(sessions) == 0 {
		return 0, 0, 0
	}

	totalSeconds := 0.0

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
		totalSeconds += dur
	}

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
