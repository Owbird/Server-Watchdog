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
	sshAttempts = []SSHAttempt{}
	ipToOrigin  = make(IPOrigins)
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

		if host != "" && !slices.Contains(whitelistedIps, host) {
			if _, exists := ipToOrigin[host]; !exists {
				stat, err := GetCountryOrigin(host)
				if err == nil {
					ipToOrigin[host] = stat
				}

			}

			if !slices.Contains(liveSSHAttempts, host) {
				liveSSHAttempts = append(liveSSHAttempts, host)
			}

		}

	}

	for i := range sshAttempts {
		if !slices.Contains(liveSSHAttempts, sshAttempts[i].IP) {
			sshAttempts[i].Status = "NIL"
		}
	}

	for _, host := range liveSSHAttempts {
		exists := false
		for i := range sshAttempts {
			if sshAttempts[i].IP == host {
				exists = true
				if sshAttempts[i].Status != "LIVE" {
					sshAttempts[i].Count += 1
				}
				sshAttempts[i].Status = "LIVE"
				break
			}
		}
		if !exists {
			sshAttempts = append(sshAttempts, SSHAttempt{
				IP:     host,
				Time:   time.Now(),
				Count:  1,
				Status: "LIVE",
			})
		}
	}

	sort.Slice(sshAttempts, func(i, j int) bool {
		if sshAttempts[i].Status == "LIVE" && sshAttempts[j].Status != "LIVE" {
			return true
		}
		if sshAttempts[j].Status == "LIVE" && sshAttempts[i].Status != "LIVE" {
			return false
		}
		return sshAttempts[i].Count > sshAttempts[j].Count
	})

	return Activities{
		WhitelistedIPs: whitelistedIps,
		Attempts:       sshAttempts,
		IPOrigins:      ipToOrigin,
	}, nil

}
