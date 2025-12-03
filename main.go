package main

import (
	"encoding/json"
	"github.com/pterm/pterm"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/bitfield/script"
)

func main() {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	logger.Info("Server Watchdog")

	whitelistedIps := []string{}
	unknownIps := []string{}

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

		if !slices.Contains(whitelistedIps, host) {
			unknownIps = append(unknownIps, host)
		}
	}

	if len(unknownIps) > 0 {

		logger.Error("Unknown ", logger.Args("IPs", strings.Join(unknownIps, ",")))
	}

}
