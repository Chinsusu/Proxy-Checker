package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"ip-proxy-checker/internal/models"
	"ip-proxy-checker/internal/proxy"
	"net/http"
	"net/url"
)

type IPQualityAPIResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	FraudScore   int    `json:"fraud_score"`
	CountryCode  string `json:"country_code"`
	Region       string `json:"region"`
	City         string `json:"city"`
	ISP          string `json:"ISP"`
	ASN          int    `json:"ASN"`
	Organization string `json:"organization"`
	Proxy        bool   `json:"proxy"`
	VPN          bool   `json:"vpn"`
	TOR          bool   `json:"tor"`
	ActiveVPN    bool   `json:"active_vpn"`
	ActiveTOR    bool   `json:"active_tor"`
}

func CheckIPQualityAPI(apiKey string, ip string, proxyClient *proxy.ProxyClient) (*models.IPQualityResult, error) {
	baseURL := fmt.Sprintf("https://ipqualityscore.com/api/json/ip/%s/%s", apiKey, ip)

	// Add parameters
	params := url.Values{}
	params.Add("strictness", "0")
	params.Add("allow_public_access_points", "true")
	params.Add("lighter_penalties", "true")
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	httpClient := http.DefaultClient
	if proxyClient != nil {
		httpClient = proxyClient.HTTPClient
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "IPQualityScore-Go-Client/1.0")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var apiResp IPQualityAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API failed: %s", apiResp.Message)
	}

	result := &models.IPQualityResult{
		IP:           ip,
		Status:       "Live",
		Country:      apiResp.CountryCode,
		City:         apiResp.City,
		Region:       apiResp.Region,
		ISP:          apiResp.ISP,
		Organization: apiResp.Organization,
		FraudScore:   fmt.Sprintf("%d", apiResp.FraudScore),
		VPN:          apiResp.VPN || apiResp.ActiveVPN,
		Proxy:        apiResp.Proxy,
		Error:        "(Source: Official API)",
	}

	return result, nil
}
