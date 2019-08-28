package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type getIpInfoRes struct {
	IP string `json:"ip"`
}

const (
	DEFAULT_CHECK_IP_URL = "https://myip.geekjoys.com"
)

func MyIP(url string) (string, error) {
	if url == "" {
		url = DEFAULT_CHECK_IP_URL
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res := getIpInfoRes{}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	return res.IP, nil
}
