package main

import (
	"log"
	"slices"
	"strconv"
	"time"

	"github.com/pterm/pterm"
)

func main() {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo)
	logger.Info("Server Watchdog Starting")

	area, _ := pterm.DefaultArea.WithCenter().Start()

	for range time.Tick(time.Second * 1) {

		activities, err := GetActivities()
		if err != nil {
			log.Fatalln(err)
		}

		header := pterm.DefaultHeader.
			Sprint("🛡️  SERVER WATCHDOG — LIVE STATUS")

		whitelistList := ""
		if len(activities.WhitelistedIPs) == 0 {
			whitelistList = pterm.Warning.Sprint("No whitelisted IPs")
		} else {
			for _, ip := range activities.WhitelistedIPs {
				whitelistList += pterm.Success.Sprint("✔ ", ip) + "\n"
			}
		}

		whitelistBox := pterm.DefaultBox.
			WithTitle("WHITELIST").
			Sprint(whitelistList)

		timestamp := pterm.FgLightBlue.Sprintf("Last Updated: %s",
			time.Now().Format("15:04:05"))

		tableData := pterm.TableData{
			{"#", "IP", "Last seen", "Country"},
		}

		for idx, ip := range activities.LiveAttempts {
			country := "N/A"
			if origin, ok := activities.IPOrigins[ip]; ok {
				country = origin.Country

			}

			tableData = append(tableData, []string{strconv.Itoa(idx + 1), ip, "NOW", country})
		}

		for _, attempt := range activities.PastAttempts {
			if !slices.Contains(activities.LiveAttempts, attempt.IP) {
				count := len(tableData)

				country := "N/A"
				if origin, ok := activities.IPOrigins[attempt.IP]; ok {

					country = origin.Country
				}

				tableData = append(tableData, []string{strconv.Itoa(count), attempt.IP, attempt.Time.Format("15:04:05"), country})
			}
		}

		table, err := pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Srender()

		if err != nil {
			log.Fatalln(err)
		}

		ui := header + "\n" +
			timestamp + "\n\n" +
			whitelistBox + "\n\n" +
			table

		area.Update(ui)

	}
}
