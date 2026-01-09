package checker

import (
	"fmt"
	"ip-proxy-checker/internal/models"
	"ip-proxy-checker/internal/proxy"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CheckAbuseIPDB(ip string, proxyClient *proxy.ProxyClient) (*models.IPQualityResult, error) {
	url := fmt.Sprintf("https://www.abuseipdb.com/check/%s", ip)

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

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("abuseipdb bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &models.IPQualityResult{
		IP:     ip,
		Status: "Live",
	}

	// Extract Abuse Confidence Score
	// AbuseIPDB shows it in a bold tag or specific div
	doc.Find(".well b, .abuse-score").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.Contains(text, "%") && result.FraudScore == "" {
			result.FraudScore = text
		}
	})

	// ISP/Org
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		label := strings.ToLower(strings.TrimSpace(s.Find("th").Text()))
		val := strings.TrimSpace(s.Find("td").Text())
		if strings.Contains(label, "isp") {
			result.ISP = val
		} else if strings.Contains(label, "domain") || strings.Contains(label, "organization") {
			result.Organization = val
		}
	})

	if result.FraudScore != "" {
		result.Error = "(Source: AbuseIPDB)"
	}

	return result, nil
}
