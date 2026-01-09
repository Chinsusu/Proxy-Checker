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
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err == nil {
		defer resp.Body.Close()
		var apiResp IPWhoResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err == nil && apiResp.Success {
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
	}

	// Fallback to ip-api.com
	return CheckIPApi(ip)
}

type IPApiResponse struct {
	Query       string `json:"query"`
	Status      string `json:"status"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	AS          string `json:"as"`
	Timezone    string `json:"timezone"`
}

func CheckIPApi(ip string) (*models.WhoisResult, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip) // Free tier uses HTTP
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return &models.WhoisResult{IP: ip, Status: "failed", Error: err.Error()}, err
	}
	defer resp.Body.Close()

	var apiResp IPApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return &models.WhoisResult{IP: ip, Status: "failed", Error: err.Error()}, err
	}

	if apiResp.Status != "success" {
		return &models.WhoisResult{IP: ip, Status: "failed", Error: "IP-API failed"}, nil
	}

	return &models.WhoisResult{
		IP:          apiResp.Query,
		Country:     apiResp.Country,
		CountryCode: apiResp.CountryCode,
		Region:      apiResp.RegionName,
		City:        apiResp.City,
		ISP:         apiResp.ISP,
		ASN:         apiResp.AS,
		Timezone:    apiResp.Timezone,
		Status:      "success",
	}, nil
}
