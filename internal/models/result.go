package models

type WhoisResult struct {
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Flag        string `json:"flag"`
	ISP         string `json:"isp"`
	ASN         string `json:"asn"`
	Timezone    string `json:"timezone"`
	Status      string `json:"status"` // "success", "failed", "pending"
	Error       string `json:"error,omitempty"`
}

type IPQualityResult struct {
	IP           string `json:"ip"`
	Port         string `json:"port,omitempty"`
	Status       string `json:"status"` // "Live", "Dead"
	Country      string `json:"country"`
	City         string `json:"city"`
	Region       string `json:"region"`
	VPN          bool   `json:"vpn"`
	Proxy        bool   `json:"proxy"`
	ISP          string `json:"isp"`
	Organization string `json:"organization"`
	Error        string `json:"error,omitempty"`
}
