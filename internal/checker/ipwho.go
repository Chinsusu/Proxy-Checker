package checker

import (
	"encoding/json"
	"fmt"
	"ip-proxy-checker/internal/models"
	"net/http"
	"time"
)

type IPWhoResponse struct {
	IP          string `json:"ip"`
	Success     bool   `json:"success"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Flag        struct {
		Img string `json:"img"`
	} `json:"flag"`
	Connection struct {
		ASN int    `json:"asn"`
		Org string `json:"org"`
		ISP string `json:"isp"`
	} `json:"connection"`
	Timezone struct {
		Utc string `json:"utc"`
	} `json:"timezone"`
	Message string `json:"message"`
}

func CheckIPWho(ip string) (*models.WhoisResult, error) {
	url := fmt.Sprintf("https://ipwho.is/%s", ip)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp IPWhoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return &models.WhoisResult{
			IP:     ip,
			Status: "failed",
			Error:  apiResp.Message,
		}, nil
	}

	return &models.WhoisResult{
		IP:          apiResp.IP,
		Country:     apiResp.Country,
		CountryCode: apiResp.CountryCode,
		Region:      apiResp.Region,
		City:        apiResp.City,
		Flag:        apiResp.Flag.Img,
		ISP:         apiResp.Connection.ISP,
		ASN:         fmt.Sprintf("AS%d", apiResp.Connection.ASN),
		Timezone:    apiResp.Timezone.Utc,
		Status:      "success",
	}, nil
}
