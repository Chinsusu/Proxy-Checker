# API Documentation - Backend Bridge

The following methods are exposed from Go to the Frontend via the Wails bridge.

## Methods

### `ParseInput(input string) map[string]interface{}`
Parses raw text and returns statistics.
- **Input**: Raw string from textarea.
- **Output**: JSON object with `ips`, `proxies`, and `total`.

### `CheckWhois(ips []string) []WhoisResult`
Concurrent Whois check for a list of IPs.
- **Input**: Array of IP strings.
- **Output**: Array of `WhoisResult` objects.

### `CheckIPQuality(proxyStrings []string) []IPQualityResult`
Concurrent IPQualityScore check using proxies.
- **Input**: Array of proxy strings (`IP:Port`).
- **Output**: Array of `IPQualityResult` objects.

## Data Structures

### `WhoisResult`
- `ip`: string
- `country`: string
- `city`: string
- `isp`: string
- `status`: "success" | "failed"

### `IPQualityResult`
- `ip`: string
- `status`: "Live" | "Dead"
- `vpn`: boolean
- `proxy`: boolean
