package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetCountryOrigin(ip string) (IPStat, error) {

	url := fmt.Sprintf("http://ip-api.com/json/%v", ip)

	res, err := http.Get(url)

	if err != nil {
		return IPStat{}, err
	}

	defer res.Body.Close()

	stat := IPStat{}

	json.NewDecoder(res.Body).Decode(&stat)

	return stat, nil

}
