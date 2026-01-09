# Proxy Information Checker

A professional desktop application built with Go and Wails v2 to check IP and Proxy information from multiple sources.

## Features
- **IPWhois Checker**: Detailed IP information (Country, City, ISP, ASN, Timezone).
- **IPQualityScore Scraper**: Advanced proxy/VPN detection.
- **Concurrent Processing**: High-performance worker pool.
- **Modern UI**: Sleek, glassmorphic design using React.
- **Format Support**: Supports simple IP lists and proxy formats (IP:Port:User:Pass).

## Installation
1. Install [Go](https://go.dev/) (1.21+).
2. Install [Node.js](https://nodejs.org/).
3. Install [Wails](https://wails.io/): `go install github.com/wailsapp/wails/v2/cmd/wails@latest`.

## Usage
1. Clone the repository.
2. Run `wails dev` for development mode.
3. Run `wails build` to create a production executable.

## Technical Details
- **Backend**: Go with `wails v2`, `sqlite` for caching, `yaml` for config.
- **Frontend**: React with Vanilla CSS for high performance and premium styling.
