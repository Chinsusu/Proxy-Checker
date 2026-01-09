package parser

import (
	"ip-proxy-checker/internal/models"
	"net"
	"strings"
)

func ParseIPList(input string) []models.IPInput {
	lines := strings.Split(input, "\n")
	var results []models.IPInput
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// If it contains a colon, it might be a proxy or IPv6
		// But simple format is just IP
		if net.ParseIP(line) != nil {
			results = append(results, models.IPInput{
				IP:   line,
				Type: "simple",
			})
		}
	}
	return results
}

func GetIPType(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "invalid"
	}
	if parsedIP.To4() != nil {
		return "v4"
	}
	return "v6"
}
