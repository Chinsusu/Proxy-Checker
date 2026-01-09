package checker

import (
	"fmt"
	"ip-proxy-checker/internal/models"
	"ip-proxy-checker/internal/proxy"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CheckIPQuality(ip string, proxyClient *proxy.ProxyClient) (*models.IPQualityResult, error) {
	url := fmt.Sprintf("https://www.ipqualityscore.com/free-ip-lookup-proxy-vpn-test/lookup/%s", ip)

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

	// Simplified headers to match reference project's success pattern
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "https://www.google.com/")

	resp, err := httpClient.Do(req)
	if err != nil {
		return &models.IPQualityResult{
			IP:     ip,
			Status: "Dead",
			Error:  err.Error(),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		// FALLBACK TO SCAMALYTICS
		fmt.Printf("[Fallback] IPQualityScore 403 Forbidden for IP: %s, trying Scamalytics...\n", ip)
		return CheckScamalytics(ip, proxyClient)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &models.IPQualityResult{
		IP:     ip,
		Status: "Live",
	}

	// Targeted extraction from the table structure
	table := doc.Find("table").First()
	if table.Length() > 0 {
		rows := table.Find("tr")

		// Row 2 contains Country, City, Region, VPN, Proxy
		if rows.Length() >= 2 {
			row2 := rows.Eq(1).Find("td")
			if row2.Length() >= 5 {
				result.Country = cleanValue(row2.Eq(0).Text())
				result.City = cleanValue(row2.Eq(1).Text())
				result.Region = cleanValue(row2.Eq(2).Text())
				result.VPN = strings.EqualFold(cleanValue(row2.Eq(3).Text()), "Yes")
				result.Proxy = strings.EqualFold(cleanValue(row2.Eq(4).Text()), "Yes")
			}
		}

		// Row 4 contains ISP, Organization, Hostname, ASN, Tor
		if rows.Length() >= 4 {
			row4 := rows.Eq(3).Find("td")
			if row4.Length() >= 5 {
				result.ISP = cleanValue(row4.Eq(0).Text())
				result.Organization = cleanValue(row4.Eq(1).Text())
				// We could add Hostname/ASN to models if needed, for now we match existing fields
			}
		}
	}

	// Fraud Score: Targeted gauge center text
	fraudScoreText := cleanValue(doc.Find(".grid-overlap.text-5xl.bold.text-center").First().Text())
	if fraudScoreText != "" {
		result.FraudScore = fraudScoreText
	}

	// Final Fallback: Label-based search if fields are still empty
	if result.Country == "" || result.Country == "N/A" {
		doc.Find("td, span, div").Each(func(i int, s *goquery.Selection) {
			txt := strings.TrimSpace(s.Text())
			if strings.HasPrefix(txt, "Country:") {
				result.Country = strings.TrimSpace(strings.TrimPrefix(txt, "Country:"))
			}
			if result.ISP == "" && strings.HasPrefix(txt, "ISP:") {
				result.ISP = strings.TrimSpace(strings.TrimPrefix(txt, "ISP:"))
			}
		})
	}

	return result, nil
}

func cleanValue(s string) string {
	f := strings.Fields(s)
	return strings.Join(f, " ")
}
