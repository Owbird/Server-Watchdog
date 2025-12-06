package main

import (
	"encoding/json"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/bitfield/script"
)

var (
	pastSSHAttempts = []PastSSHAttempt{}
	ipToOrigin      = make(IPOrigins)
)

func GetActivities() (Activities, error) {
	whitelistedIps := []string{}

	if _, err := os.Stat("whitelist.json"); err != nil {
		os.WriteFile("whitelist.json", []byte("[]"), 0755)
	} else {
		data, err := os.ReadFile("whitelist.json")
		if err != nil {
			return Activities{}, err
		}
		json.Unmarshal(data, &whitelistedIps)
	}

	liveSSHAttempts := []string{}

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

			found := false

			for _, atmpt := range pastSSHAttempts {

				if atmpt.IP == host {
					found = true
					break
				}
			}

			if !found {
				pastSSHAttempts = append(pastSSHAttempts, PastSSHAttempt{
					IP:   host,
					Time: time.Now(),
				})
			}

		}

	}

	return Activities{
		WhitelistedIPs: whitelistedIps,
		LiveAttempts:   liveSSHAttempts,
		PastAttempts:   pastSSHAttempts,
		IPOrigins:      ipToOrigin,
	}, nil

}
