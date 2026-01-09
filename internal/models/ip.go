package models

type IPAddress struct {
	IP   string `json:"ip"`
	Type string `json:"type"` // "v4" or "v6"
}

type IPInput struct {
	IP   string `json:"ip"`
	Type string `json:"type"` // "simple" or "proxy"
}
