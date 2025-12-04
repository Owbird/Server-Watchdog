package main

import (
	"encoding/json"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/pterm/pterm"
)

func main() {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo)
	logger.Info("Server Watchdog Starting")

	whitelistedIps := []string{}
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

		unknownSShAttempts := []string{}

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

			if host != "" &&
				!slices.Contains(whitelistedIps, host) &&
				!slices.Contains(unknownSShAttempts, host) {
				unknownSShAttempts = append(unknownSShAttempts, host)
			}
		}


		header := pterm.DefaultHeader.
			Sprint("🛡️  SERVER WATCHDOG — LIVE STATUS")

		unknownList := ""
		if len(unknownSShAttempts) == 0 {
			unknownList = pterm.Success.Sprint("No unknown attempts detected")
		} else {
			for _, ip := range unknownSShAttempts {
				unknownList += pterm.Error.Sprint("❌ ", ip) + "\n"
			}
		}

		unknownBox := pterm.DefaultBox.
			WithTitle("SSH MONITOR").
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
			timestamp + "\n"


		area.Update(ui)

	}
}
