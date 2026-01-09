package parser

import (
	"testing"
)

func TestParseIPList(t *testing.T) {
	input := "1.1.1.1\n8.8.8.8\ninvalid\n2606:4700:4700::1111"
	ips := ParseIPList(input)

	if len(ips) != 3 {
		t.Errorf("Expected 3 valid IPs, got %d", len(ips))
	}
}

func TestParseProxyList(t *testing.T) {
	input := "1.2.3.4:8080\n5.6.7.8:9090:user:pass\ninvalid"
	proxies := ParseProxyList(input)

	if len(proxies) != 2 {
		t.Errorf("Expected 2 valid proxies, got %d", len(proxies))
	}

	if proxies[1].Username != "user" || proxies[1].Password != "pass" {
		t.Errorf("Proxy auth failed to parse")
	}
}
