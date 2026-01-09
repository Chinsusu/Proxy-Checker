package main

import (
	"encoding/json"
	"io"
	"ip-proxy-checker/internal/checker"
	"ip-proxy-checker/internal/models"
	"ip-proxy-checker/internal/parser"
	"ip-proxy-checker/internal/proxy"
	"ip-proxy-checker/internal/storage"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	config *storage.Config
	cache  *storage.Cache
	logger zerolog.Logger
}

func NewApp() *App {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &App{logger: logger}
}

func (a *App) Init() error {
	cfg, err := storage.LoadConfig("config.yaml")
	if err != nil {
		a.logger.Warn().Err(err).Msg("Failed to load config, using defaults")
	}
	a.config = cfg

	cache, err := storage.NewCache(a.config.Storage.DBPath)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to initialize cache")
		return err
	}
	a.cache = cache
	return nil
}

func (a *App) HandleParseInput(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Input string `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ips := parser.ParseIPList(body.Input)
	proxies := parser.ParseProxyList(body.Input)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ips":     ips,
		"proxies": proxies,
		"total":   len(ips) + len(proxies),
	})
}

func (a *App) HandleCheckWhois(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IPs []string `json:"ips"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.logger.Info().Int("count", len(body.IPs)).Msg("Starting Whois check")
	results := make([]models.WhoisResult, 0)
	wp := checker.NewWorkerPool(a.config.Worker.PoolSize)

	wp.Start(func(job checker.Job) interface{} {
		ip := job.Data.(string)
		res, err := checker.CheckIPWho(ip)
		if err != nil {
			a.logger.Error().Err(err).Str("ip", ip).Msg("Whois check failed")
			return models.WhoisResult{IP: ip, Status: "failed", Error: err.Error()}
		}
		if res.Status == "failed" {
			a.logger.Warn().Str("ip", ip).Str("error", res.Error).Msg("Whois API returned failure")
		}
		return *res
	})

	for i, ip := range body.IPs {
		wp.AddJob(checker.Job{ID: i, Type: "whois", Data: ip})
	}

	for i := 0; i < len(body.IPs); i++ {
		res := <-wp.Results()
		results = append(results, res.(models.WhoisResult))
	}
	wp.Stop()

	json.NewEncoder(w).Encode(results)
}

func (a *App) HandleCheckIPQuality(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Proxies []string `json:"proxies"`
		APIKey  string   `json:"api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.logger.Debug().
		Int("proxies_count", len(body.Proxies)).
		Str("api_key_provided", func() string {
			if body.APIKey == "" {
				return "NO"
			}
			return "YES"
		}()).
		Msg("Received Quality Check Request")

	// Update config if APIKey is provided in request
	effectiveAPIKey := a.config.API.IPQuality.APIKey
	if body.APIKey != "" {
		effectiveAPIKey = body.APIKey
		a.logger.Debug().Msg("Using API Key provided in request")
	}

	a.logger.Info().Int("count", len(body.Proxies)).Msg("Starting IPQuality check")
	results := make([]models.IPQualityResult, 0)
	wp := checker.NewWorkerPool(a.config.Worker.PoolSize)

	wp.Start(func(job checker.Job) interface{} {
		proxyStr := job.Data.(string)

		// Extract Host and Port for default display
		host := proxyStr
		port := ""
		if idx := strings.Index(host, ":"); idx != -1 {
			port = host[idx+1:]
			host = host[:idx]
			if idx2 := strings.Index(port, ":"); idx2 != -1 {
				port = port[:idx2]
			}
		}

		ua := proxy.GetRandomUserAgent()
		// NewProxyClient now handles protocol normalization internally
		client, err := proxy.NewProxyClient(proxyStr, ua, a.config.Proxy.ConnectionTimeout)
		if err != nil {
			a.logger.Error().Err(err).Str("proxy", proxyStr).Msg("Failed to create proxy client")
			return models.IPQualityResult{IP: host, Port: port, Status: "Dead", Error: "Invalid proxy format"}
		}

		// Step 0: TCP Pre-Check (Verify port reachability)
		a.logger.Info().Str("proxy", proxyStr).Msg(">>> STEP 0: Verifying TCP Port Reachability")
		if err := client.RawTCPCheck(); err != nil {
			a.logger.Warn().Err(err).Str("proxy", proxyStr).Msg("Proxy Port Unreachable")
			return models.IPQualityResult{IP: host, Port: port, Status: "Dead", Error: "TCP unreachable: " + err.Error()}
		}
		a.logger.Info().Str("proxy", proxyStr).Msg("TCP Port is OPEN")

		// Step 1: Connectivity Check & Exit IP Detection
		// Switch to Amazon CheckIP for better reliability and to get our actual IP
		testTarget := "http://checkip.amazonaws.com"
		a.logger.Info().Str("proxy", proxyStr).Msg(">>> STEP 1: Testing Connectivity & Detecting Exit IP")
		testResp, err := client.HTTPClient.Get(testTarget)

		// Protocol Fallback: If no protocol was specified and HTTP failed, try SOCKS5
		if err != nil && !strings.Contains(proxyStr, "://") {
			a.logger.Info().Str("proxy", proxyStr).Msg("HTTP check failed. Trying SOCKS5 fallback...")
			s5Client, s5Err := proxy.NewProxyClient("socks5://"+proxyStr, ua, a.config.Proxy.ConnectionTimeout)
			if s5Err == nil {
				a.logger.Info().Str("proxy", proxyStr).Msg(">>> STEP 2: Testing SOCKS5 Connectivity")
				testResp, s5Err = s5Client.HTTPClient.Get(testTarget)
				if s5Err == nil {
					a.logger.Info().Str("proxy", proxyStr).Msg("SOCKS5 connectivity verified successfully!")
					client = s5Client // Switch to SOCKS5 client for subsequent checks
					err = nil
				} else {
					a.logger.Warn().Str("proxy", proxyStr).Err(s5Err).Msg("SOCKS5 fallback also failed.")
				}
			}
		}

		if err != nil {
			a.logger.Warn().Err(err).Str("proxy", proxyStr).Msg("Proxy is DEAD - Protocol/Auth failed.")
			return models.IPQualityResult{IP: host, Port: port, Status: "Dead", Error: "Protocol failed: " + err.Error()}
		}

		// Read the Exit IP from the response
		exitIPBytes, _ := io.ReadAll(io.LimitReader(testResp.Body, 1024))
		testResp.Body.Close()
		exitIP := strings.TrimSpace(string(exitIPBytes))

		if exitIP == "" {
			a.logger.Warn().Str("proxy", proxyStr).Msg("Could not detect Exit IP")
			return models.IPQualityResult{IP: host, Port: port, Status: "Dead", Error: "Could not detect Exit IP"}
		}
		a.logger.Info().Str("proxy", proxyStr).Str("exit_ip", exitIP).Msg("Proxy is LIVE")

		// Step 2: Quality Check (Use the ACTUAL Exit IP)
		var res *models.IPQualityResult
		// err is already declared in the outer scope of the closure

		// Try Official API first if API key is set
		if effectiveAPIKey != "" {
			maskedKey := "set"
			if len(effectiveAPIKey) > 8 {
				maskedKey = effectiveAPIKey[:4] + "****" + effectiveAPIKey[len(effectiveAPIKey)-4:]
			}
			a.logger.Info().Str("exit_ip", exitIP).Str("api_key", maskedKey).Msg("Using Official IPQuality API...")
			res, err = checker.CheckIPQualityAPI(effectiveAPIKey, exitIP, client)
		} else {
			a.logger.Info().Str("exit_ip", exitIP).Msg("No API Key set, using scraping method...")
			// No API key, use scraping
			res, err = checker.CheckIPQuality(exitIP, client)
		}

		if err != nil {
			// HYBRID FALLBACK: If proxy check is blocked (403) or API failed, try checking directly from Local IP
			if strings.Contains(err.Error(), "403") || effectiveAPIKey != "" {
				a.logger.Info().Str("exit_ip", exitIP).Str("proxy", proxyStr).Msg("[Fallback] Check failed. Trying Direct Check...")

				// Try API directly first
				if effectiveAPIKey != "" {
					res, err = checker.CheckIPQualityAPI(effectiveAPIKey, exitIP, nil)
				}

				// If Direct API failed or wasn't used, try Scraping directly
				if err != nil || effectiveAPIKey == "" {
					res, err = checker.CheckIPQuality(exitIP, nil) // passing nil for proxyClient means Direct Check
				}

				if err != nil || (res != nil && res.FraudScore == "") {
					a.logger.Info().Str("exit_ip", exitIP).Msg("[Fallback] IPQualityScore failed. Trying Scamalytics...")
					res, err = checker.CheckScamalytics(exitIP, nil)
				}

				if err != nil || (res != nil && res.FraudScore == "") {
					a.logger.Info().Str("exit_ip", exitIP).Msg("[Fallback] Scamalytics failed. Trying AbuseIPDB...")
					res, err = checker.CheckAbuseIPDB(exitIP, nil)
				}

				if err == nil {
					if res.Error == "" {
						res.Error = "(Source: Direct Check)"
					}
				}
			}

			if err != nil {
				a.logger.Error().Err(err).Str("exit_ip", exitIP).Str("proxy", proxyStr).Msg("IPQuality check failed even with fallback")
				// Even if IPQuality fails, it's still 'Live' because step 1 passed
				return models.IPQualityResult{IP: exitIP, Port: port, Status: "Live", Error: "Quality info failed: " + err.Error()}
			}
		}
		res.Port = port
		res.Status = "Live" // Ensure status is Live if we reach here
		return *res
	})

	for i, p := range body.Proxies {
		wp.AddJob(checker.Job{ID: i, Type: "ipquality", Data: p})
	}

	for i := 0; i < len(body.Proxies); i++ {
		res := <-wp.Results()
		results = append(results, res.(models.IPQualityResult))
	}
	wp.Stop()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (a *App) HandleSetAPIKey(w http.ResponseWriter, r *http.Request) {
	var body struct {
		APIKey string `json:"api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.config.API.IPQuality.APIKey = strings.TrimSpace(body.APIKey)
	a.logger.Info().Msg("IPQuality API Key updated from frontend")

	w.WriteHeader(http.StatusOK)
}
