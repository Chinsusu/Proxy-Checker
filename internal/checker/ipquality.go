package checker

import (
	"fmt"
	"io"
	"ip-proxy-checker/internal/models"
	"ip-proxy-checker/internal/proxy"
	"net/http"
	"regexp"
	"strings"
)

func CheckIPQuality(ip string, proxyClient *proxy.ProxyClient) (*models.IPQualityResult, error) {
	url := fmt.Sprintf("https://www.ipqualityscore.com/free-ip-lookup-proxy-vpn-test/lookup/%s", ip)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", proxyClient.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")

	resp, err := proxyClient.HTTPClient.Do(req)
	if err != nil {
		return &models.IPQualityResult{
			IP:     ip,
			Status: "Dead",
			Error:  err.Error(),
		}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	html := string(body)
	result := &models.IPQualityResult{
		IP:     ip,
		Status: "Live",
	}

	// Simple regex based scraper for demo purposes
	// In production, use a more robust HTML parser like goquery
	result.Country = extractField(html, `Country:</td>\s*<td>([^<]+)`)
	result.City = extractField(html, `City:</td>\s*<td>([^<]+)`)
	result.Region = extractField(html, `Region:</td>\s*<td>([^<]+)`)
	result.ISP = extractField(html, `ISP:</td>\s*<td>([^<]+)`)
	result.Organization = extractField(html, `Organization:</td>\s*<td>([^<]+)`)

	result.VPN = strings.Contains(strings.ToLower(html), "vpn detected: yes")
	result.Proxy = strings.Contains(strings.ToLower(html), "proxy detected: yes")

	return result, nil
}

func extractField(html, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return "N/A"
}
