package checker

import (
	"fmt"
	"ip-proxy-checker/internal/models"
	"ip-proxy-checker/internal/proxy"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CheckScamalytics(ip string, proxyClient *proxy.ProxyClient) (*models.IPQualityResult, error) {
	url := fmt.Sprintf("https://scamalytics.com/ip/%s", ip)

	httpClient := http.DefaultClient
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"
	if proxyClient != nil {
		httpClient = proxyClient.HTTPClient
		userAgent = proxyClient.UserAgent
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Standardized headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://scamalytics.com/")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scamalytics bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &models.IPQualityResult{
		IP:     ip,
		Status: "Live",
	}

	// Extract Fraud Score (it's often in a div with a numeric score)
	// Based on Scamalytics structure: typically <div class="score"> or similar
	doc.Find(".score, .risk-score, div[style*='score']").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" && result.FraudScore == "" {
			result.FraudScore = text
		}
	})

	// Alternative extraction for Scamalytics:
	// The score is often inside a table or a specific div
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		line := strings.TrimSpace(s.Text())
		if strings.Contains(line, "Fraud Score") {
			val := strings.TrimSpace(s.Find("td").Last().Text())
			if val != "" {
				result.FraudScore = val
			}
		}
	})

	// ISP/Org
	doc.Find("th").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Text())
		val := strings.TrimSpace(s.Next().Text())
		switch label {
		case "ISP":
			result.ISP = val
		case "Organization":
			result.Organization = val
		case "Country":
			result.Country = val
		}
	})

	// Check for Proxy/VPN/Tor
	bodyText := strings.ToLower(doc.Text())
	result.Proxy = strings.Contains(bodyText, "proxy: yes") || strings.Contains(bodyText, "is a proxy")
	result.VPN = strings.Contains(bodyText, "vpn: yes") || strings.Contains(bodyText, "is a vpn")

	// If score found, add a note that this is from Scamalytics
	if result.FraudScore != "" {
		result.Error = "(Source: Scamalytics)"
	}

	return result, nil
}
