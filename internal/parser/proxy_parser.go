package parser

import (
	"ip-proxy-checker/internal/models"
	"strings"
)

// ParseProxyList parses strings in format IP:PORT or IP:PORT:USER:PASS
func ParseProxyList(input string) []models.ProxyInput {
	lines := strings.Split(input, "\n")
	var results []models.ProxyInput
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			proxy := models.ProxyInput{
				Host: parts[0],
				Port: parts[1],
			}
			if len(parts) >= 4 {
				proxy.Username = parts[2]
				proxy.Password = parts[3]
			}
			results = append(results, proxy)
		}
	}
	return results
}
