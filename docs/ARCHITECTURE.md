# Architecture Documentation - Proxy Checker (Web Server)

## System Overview
The Proxy Checker is a web-based application designed for high-concurrency IP and proxy verification. It features a Go-based REST API and a React-based SPA (Single Page Application) served directly by the Go backend.

## Components
### 1. Backend (Go + Chi)
- **REST API**: Exposes endpoints for parsing input and triggering concurrent checks.
- **Worker Pool**: Manages parallel execution of checker jobs.
- **Static File Server**: Serves the bundled React frontend from an embedded filesystem.
- **Storage**: SQLite cache for persisting results.

### 2. Frontend (React)
- **API Service**: Uses `fetch` to communicate with the Go backend.
- **State Management**: React hooks for handling data and UI state.
- **Responsive UI**: Glassmorphic design that works across standard web browsers.

## Data Flow
1. The user accesses the web server via a browser.
2. The Go server sends the React application to the client.
3. User interactions trigger `fetch` requests to `/api/*` endpoints.
4. The backend processes requests using its worker pool and returns JSON responses.
