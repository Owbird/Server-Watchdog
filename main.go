package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	"github.com/xeonx/timeago"
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
			{"#", "IP", "Status", "Last seen", "Country", "Sessions"},
		}

		for idx, attempt := range activities.Attempts {

			country := "N/A"
			if origin, ok := activities.IPOrigins[attempt.IP]; ok {

				country = origin.Country
			}

			fmtedTime := ""
			if attempt.Status == LIVE {

				fmtedTime = "NOW"
			} else {

				lastSession := attempt.Sessions[len(attempt.Sessions)-1]
				fmtedTime = timeago.English.Format(lastSession.End)

			}

			r, g, b := GetWarmth(attempt.Sessions)

			sessionsString := fmt.Sprintf("%v attempts (%v)", strconv.Itoa(len(attempt.Sessions)), FormatDuration(GetTotalDuration(attempt.Sessions)))

			tableRow := RGBify(r, g, b, []string{
				strconv.Itoa(idx + 1), attempt.IP, attempt.Status, fmtedTime, country, sessionsString,
			})
			tableData = append(tableData, tableRow)
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
