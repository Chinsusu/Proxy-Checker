# Architecture Documentation - Proxy Checker

## System Overview
The Proxy Checker is a desktop application designed for high-concurrency IP and proxy verification. It uses a Go backend for performance and a React frontend for a rich user experience, bridged by Wails v2.

## Components
### 1. Backend (Go)
- **Parser**: Sanitizes and categorizes user input into IP or Proxy models.
- **Checker**: Implements API clients (IPWho.is) and web scrapers (IPQualityScore).
- **Worker Pool**: Manages concurrent tasks to ensure responsiveness and high throughput.
- **Proxy Client**: Custom HTTP client with support for authentication and User-Agent rotation.
- **Storage**: SQLite-based caching layer to reduce API calls and improve performance.

### 2. Frontend (React)
- **State Management**: React hooks for local state and results.
- **UI Components**: Modular tabs for Home, Whois, and Quality analysis.
- **Styling**: Vanilla CSS utilizing modern properties (glassmorphism, CSS variables).

## Data Flow
1. User enters raw text in the Home tab.
2. The frontend calls `ParseInput` via the Wails bridge.
3. The backend returns parsed statistics.
4. When a check is triggered, the backend spins up workers.
5. Workers process jobs and return results to the frontend.
