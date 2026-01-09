package main

import (
	"encoding/json"
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
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.logger.Info().Int("count", len(body.Proxies)).Msg("Starting IPQuality check")
	results := make([]models.IPQualityResult, 0)
	wp := checker.NewWorkerPool(a.config.Worker.PoolSize)

	wp.Start(func(job checker.Job) interface{} {
		proxyStr := job.Data.(string)

		// Extract IP and Port for default display
		ip := proxyStr
		port := ""
		if idx := strings.Index(ip, ":"); idx != -1 {
			port = ip[idx+1:]
			ip = ip[:idx]
			if idx2 := strings.Index(port, ":"); idx2 != -1 {
				port = port[:idx2]
			}
		}

		ua := proxy.GetRandomUserAgent()
		// NewProxyClient now handles protocol normalization internally
		client, err := proxy.NewProxyClient(proxyStr, ua, a.config.Proxy.ConnectionTimeout)
		if err != nil {
			a.logger.Error().Err(err).Str("proxy", proxyStr).Msg("Failed to create proxy client")
			return models.IPQualityResult{IP: ip, Port: port, Status: "Dead", Error: "Invalid proxy format"}
		}

		// Step 0: TCP Pre-Check (Verify port reachability)
		a.logger.Info().Str("proxy", proxyStr).Msg(">>> STEP 0: Verifying TCP Port Reachability")
		if err := client.RawTCPCheck(); err != nil {
			a.logger.Warn().Err(err).Str("proxy", proxyStr).Msg("Proxy Port Unreachable")
			return models.IPQualityResult{IP: ip, Port: port, Status: "Dead", Error: "TCP unreachable: " + err.Error()}
		}
		a.logger.Info().Str("proxy", proxyStr).Msg("TCP Port is OPEN")

		// Step 1: Connectivity Check (Live check)
		// Switch to Amazon CheckIP for better reliability
		testTarget := "http://checkip.amazonaws.com"
		a.logger.Info().Str("proxy", proxyStr).Msg(">>> STEP 1: Testing HTTP Connectivity")
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
					testResp.Body.Close()
					client = s5Client // Switch to SOCKS5 client for subsequent checks
					err = nil
				} else {
					a.logger.Warn().Str("proxy", proxyStr).Err(s5Err).Msg("SOCKS5 fallback also failed.")
				}
			}
		}

		if err != nil {
			a.logger.Warn().Err(err).Str("proxy", proxyStr).Msg("Proxy is DEAD - Protocol/Auth failed.")
			return models.IPQualityResult{IP: ip, Port: port, Status: "Dead", Error: "Protocol failed: " + err.Error()}
		}
		testResp.Body.Close()
		a.logger.Info().Str("proxy", proxyStr).Msg("Proxy is LIVE")

		// Step 2: Quality Check
		res, err := checker.CheckIPQuality(ip, client)
		if err != nil {
			a.logger.Error().Err(err).Str("ip", ip).Str("proxy", proxyStr).Msg("IPQuality check failed")
			// Even if IPQuality fails, it's still 'Live' because step 1 passed
			return models.IPQualityResult{IP: ip, Port: port, Status: "Live", Error: "Quality info failed: " + err.Error()}
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

	json.NewEncoder(w).Encode(results)
}
