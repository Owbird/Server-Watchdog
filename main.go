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

		unknownIps := []string{}

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

			if !slices.Contains(whitelistedIps, host) && !slices.Contains(unknownIps, host) {
				unknownIps = append(unknownIps, host)
			}
		}

		if len(unknownIps) > 0 {

			txt := pterm.Sprintf("Unknown IPs: %v", strings.Join(unknownIps, ","))
			area.Update(txt)

		}
	}

}
