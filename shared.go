package main

import "time"

type IPOrigins map[string]IPStat

type PastSSHAttempt struct {
	IP   string
	Time time.Time
}

type Activities struct {
	WhitelistedIPs []string
	LiveAttempts   []string
	PastAttempts   []PastSSHAttempt
	IPOrigins      IPOrigins
}

type IPStat struct {
	Query       string
	Status      string
	Country     string
	CountryCode string
	Region      string
	RegionName  string
	City        string
	Zip         string
	Lat         float64
	Lon         float64
	Timezone    string
	Isp         string
	Org         string
	As          string
}
