package main

import "time"

const (
	LIVE  = "LIVE"
	STALE = "NIL"
)

type AttemptSession struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type SSHAttempt struct {
	IP       string           `json:"ip"`
	Status   string           `json:"status"`
	Country  string           `json:"country"`
	Sessions []AttemptSession `json:"sessions"`
}

type Activities struct {
	WhitelistedIPs []string     `json:"whitelistedIPs"`
	Attempts       []SSHAttempt `json:"attempts"`
}

type IPStat struct {
	Query       string  `json:"query,omitempty"`
	Status      string  `json:"status,omitempty"`
	Country     string  `json:"country,omitempty"`
	CountryCode string  `json:"countryCode,omitempty"`
	Region      string  `json:"region,omitempty"`
	RegionName  string  `json:"regionName,omitempty"`
	City        string  `json:"city,omitempty"`
	Zip         string  `json:"zip,omitempty"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone,omitempty"`
	Isp         string  `json:"isp,omitempty"`
	Org         string  `json:"org,omitempty"`
	As          string  `json:"as,omitempty"`
}
