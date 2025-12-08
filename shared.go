package main

import "time"

type IPOrigins map[string]IPStat

type AttemptSession struct {
	Start time.Time
	End   time.Time
}

type SSHAttempt struct {
	IP       string
	Status   string
	Sessions []AttemptSession
}

type Activities struct {
	WhitelistedIPs []string
	Attempts       []SSHAttempt
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
