# REST API Documentation

The Proxy Checker exposes a RESTful API for integration.

## Base URL
`/api`

## Endpoints

### `POST /parse`
Parses raw text into structured IP and Proxy counts.
- **Request Body**: `{ "input": "string" }`
- **Response**: `{ "ips": [...], "proxies": [...], "total": 0 }`

### `POST /check/whois`
Performs concurrent Whois lookups.
- **Request Body**: `{ "ips": ["1.1.1.1", ...] }`
- **Response**: `[ { "ip": "...", "country": "...", ... }, ... ]`

### `POST /check/quality`
Performs concurrent IPQuality analysis using proxies.
- **Request Body**: `{ "proxies": ["IP:Port", ...] }`
- **Response**: `[ { "ip": "...", "status": "Live", ... }, ... ]`

## Status Codes
- `200 OK`: Request successful.
- `400 Bad Request`: Invalid JSON or missing parameters.
- `500 Internal Server Error`: Server-side processing error.
