package main

import (
	"encoding/json"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bitfield/script"
)

var (
	sshAttempts = LoadAttempts()
)

func GetActivities() (Activities, error) {
	whitelistedIps := []string{}
	liveSSHAttempts := []string{}

	if _, err := os.Stat("whitelist.json"); err != nil {
		os.WriteFile("whitelist.json", []byte("[]"), 0755)
	} else {
		data, err := os.ReadFile("whitelist.json")
		if err != nil {
			return Activities{}, err
		}
		json.Unmarshal(data, &whitelistedIps)
	}

	tnp, err := script.Exec("ss -tnp").Match(":22").String()

	if err != nil {
		return Activities{}, err
	}

	for line := range strings.SplitSeq(tnp, "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		ip, err := script.Echo(line).Column(5).String()
		if err != nil {
			continue
		}

		host := strings.Split(ip, ":")[0]

		if host != "" && !slices.Contains(whitelistedIps, host) && !strings.Contains(host, "[") {

			if !slices.Contains(liveSSHAttempts, host) {
				liveSSHAttempts = append(liveSSHAttempts, host)
			}

		}

	}

	for i := range sshAttempts {
		if !slices.Contains(liveSSHAttempts, sshAttempts[i].IP) {
			if sshAttempts[i].Status == LIVE {
				totalSessions := len(sshAttempts[i].Sessions)
				if totalSessions > 0 {
					sshAttempts[i].Sessions[totalSessions-1].End = time.Now()
				}
			}
			sshAttempts[i].Status = STALE
		}
	}

	for _, host := range liveSSHAttempts {
		exists := false
		for i := range sshAttempts {
			if sshAttempts[i].IP == host {
				exists = true
				if sshAttempts[i].Status != LIVE {
					sshAttempts[i].Sessions = append(sshAttempts[i].Sessions, AttemptSession{
						Start: time.Now(),
					})

				}
				sshAttempts[i].Status = LIVE
				break
			}
		}
		if !exists {

			country := "N/A"

			ipStat, err := GetCountryOrigin(host)
			if err == nil {
				country = ipStat.Country
			}

			sshAttempts = append(sshAttempts, SSHAttempt{
				IP: host,
				Sessions: []AttemptSession{
					{
						Start: time.Now(),
					}},
				Country: country,
				Status:  LIVE,
			})
		}
	}

	sort.Slice(sshAttempts, func(i, j int) bool {
		if sshAttempts[i].Status == LIVE && sshAttempts[j].Status != LIVE {
			return true
		}
		if sshAttempts[j].Status == LIVE && sshAttempts[i].Status != LIVE {
			return false
		}
		return len(sshAttempts[i].Sessions) > len(sshAttempts[j].Sessions)
	})

	SaveAttemts(sshAttempts)

	return Activities{
		WhitelistedIPs: whitelistedIps,
		Attempts:       sshAttempts,
	}, nil

}
