package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/pterm/pterm"
)

type PastSSHAttempt struct {
	IP   string
	Time time.Time
}

func main() {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo)
	logger.Info("Server Watchdog Starting")

	whitelistedIps := []string{}
	pastSSHAttempts := []PastSSHAttempt{}

	if _, err := os.Stat("whitelist.json"); err != nil {
		os.WriteFile("whitelist.json", []byte("[]"), 0755)
	} else {
		data, err := os.ReadFile("whitelist.json")
		if err != nil {
			log.Fatalln(err)
		}
		json.Unmarshal(data, &whitelistedIps)
	}

	area, _ := pterm.DefaultArea.WithCenter().Start()

	for range time.Tick(time.Second * 1) {

		unknownSSHAttempts := []string{}

		tnp, err := script.Exec("ss -tnp").Match(":22").String()

		if err != nil {
			log.Fatalln(err)
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
				if !slices.Contains(unknownSSHAttempts, host) {
					unknownSSHAttempts = append(unknownSSHAttempts, host)
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

		header := pterm.DefaultHeader.
			Sprint("🛡️  SERVER WATCHDOG — LIVE STATUS")

		unknownList := ""
		if len(unknownSSHAttempts) == 0 {
			unknownList = pterm.Success.Sprint("No unknown attempts detected")
		} else {
			for _, ip := range unknownSSHAttempts {
				unknownList += pterm.Red("[NOW] ", ip) + "\n"
			}

			for _, attempt := range pastSSHAttempts {
				if !slices.Contains(unknownSSHAttempts, attempt.IP) {
					unknownList += pterm.Yellow(fmt.Sprintf("[Last seen: %v] ", attempt.Time.Format("15:04:05")), attempt.IP) + "\n"
				}
			}

		}

		unknownBox := pterm.DefaultBox.
			WithTitle("SSH MONITOR (Live)").
			Sprint(unknownList)

		whitelistList := ""
		if len(whitelistedIps) == 0 {
			whitelistList = pterm.Warning.Sprint("No whitelisted IPs")
		} else {
			for _, ip := range whitelistedIps {
				whitelistList += pterm.Success.Sprint("✔ ", ip) + "\n"
			}
		}

		whitelistBox := pterm.DefaultBox.
			WithTitle("WHITELIST").
			Sprint(whitelistList)

		timestamp := pterm.FgLightBlue.Sprintf("Last Updated: %s",
			time.Now().Format("15:04:05"))

		ui := header + "\n\n" +
			unknownBox + "\n\n" +
			whitelistBox + "\n\n" +
			timestamp

		area.Update(ui)

	}
}

