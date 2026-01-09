# Proxy Information Checker (Web Server)

A high-performance web application to check IP and Proxy information using a Go backend and a React frontend.

## Features
- **REST API**: Clean HTTP endpoints for integration.
- **Embedded Frontend**: Backend binary serves the entire web app.
- **IPWhois & IPQuality**: Integrated data sources with parallel processing.
- **CSV Export**: Download results directly from the browser.

## Getting Started

### Prerequisites
- [Go](https://go.dev/) (1.21+)
- [Node.js](https://nodejs.org/) (for frontend development)

### Development
1. **Frontend**:
   ```bash
   cd frontend
   npm install
   npm run build
   ```
2. **Backend**:
   ```bash
   go run .
   ```
3. Access at `http://localhost:8080`.

### Build
Run `make build` to generate the production binary containing the embedded frontend.
