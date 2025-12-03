package main

import (
	"encoding/json"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/pterm/pterm"

	"github.com/bitfield/script"
)

func main() {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	logger.Info("Server Watchdog")

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

	logger.Info("White Listed ", logger.Args("IPs", strings.Join(whitelistedIps, ",")))

	area, _ := pterm.DefaultArea.Start()

	for range time.Tick(time.Second * 1) {

		unknownSShAttempts := []string{}

		tnp, err := script.Exec("ss -tnp").Match(":22").String()
		if err != nil {
			log.Fatalln(err)
		}

		for line := range strings.SplitSeq(tnp, "\n") {
			ip, err := script.Echo(line).Column(5).String()
			if err != nil {
				log.Fatalln(err)
			}

			host := strings.Split(ip, ":")[0]

			if !slices.Contains(whitelistedIps, host) && !slices.Contains(unknownSShAttempts, host) {
				unknownSShAttempts = append(unknownSShAttempts, host)
			}
		}

		if len(unknownSShAttempts) > 0 {

			txt := "Unknown attempts: " + pterm.Error.Sprint(strings.Join(unknownSShAttempts, ", "))
			area.Update(txt)

		}
	}

}
