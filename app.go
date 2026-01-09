package main

import (
	"context"
	"ip-proxy-checker/internal/checker"
	"ip-proxy-checker/internal/models"
	"ip-proxy-checker/internal/parser"
	"ip-proxy-checker/internal/proxy"
	"ip-proxy-checker/internal/storage"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	ctx    context.Context
	config *storage.Config
	cache  *storage.Cache
	logger zerolog.Logger
}

func NewApp() *App {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &App{logger: logger}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.logger.Info().Msg("Starting application...")

	cfg, err := storage.LoadConfig("config.yaml")
	if err != nil {
		a.logger.Warn().Err(err).Msg("Failed to load config, using defaults")
	}
	a.config = cfg

	cache, err := storage.NewCache(a.config.Storage.DBPath)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to initialize cache")
	}
	a.cache = cache
}

func (a *App) ParseInput(input string) map[string]interface{} {
	a.logger.Debug().Msg("Parsing user input")
	ips := parser.ParseIPList(input)
	proxies := parser.ParseProxyList(input)
	return map[string]interface{}{
		"ips":     ips,
		"proxies": proxies,
		"total":   len(ips) + len(proxies),
	}
}

func (a *App) CheckWhois(ips []string) []models.WhoisResult {
	a.logger.Info().Int("count", len(ips)).Msg("Starting Whois check")
	results := make([]models.WhoisResult, 0)
	wp := checker.NewWorkerPool(a.config.Worker.PoolSize)

	wp.Start(func(job checker.Job) interface{} {
		ip := job.Data.(string)
		a.logger.Debug().Str("ip", ip).Msg("Checking Whois")
		res, err := checker.CheckIPWho(ip)
		if err != nil {
			a.logger.Error().Err(err).Str("ip", ip).Msg("Whois check failed")
			return models.WhoisResult{IP: ip, Status: "failed", Error: err.Error()}
		}
		return *res
	})

	for i, ip := range ips {
		wp.AddJob(checker.Job{ID: i, Type: "whois", Data: ip})
	}

	for i := 0; i < len(ips); i++ {
		res := <-wp.Results()
		results = append(results, res.(models.WhoisResult))
	}
	wp.Stop()

	a.logger.Info().Int("results", len(results)).Msg("Whois check completed")
	return results
}

func (a *App) CheckIPQuality(proxyStrings []string) []models.IPQualityResult {
	a.logger.Info().Int("count", len(proxyStrings)).Msg("Starting IPQuality check")
	results := make([]models.IPQualityResult, 0)
	wp := checker.NewWorkerPool(a.config.Worker.PoolSize)

	wp.Start(func(job checker.Job) interface{} {
		proxyStr := job.Data.(string)
		ua := proxy.GetRandomUserAgent()
		client, err := proxy.NewProxyClient("http://"+proxyStr, ua, a.config.Proxy.ConnectionTimeout)
		if err != nil {
			return models.IPQualityResult{IP: proxyStr, Status: "Dead", Error: err.Error()}
		}

		a.logger.Debug().Str("proxy", proxyStr).Msg("Checking IPQuality")
		res, err := checker.CheckIPQuality(proxyStr, client)
		if err != nil {
			a.logger.Error().Err(err).Str("proxy", proxyStr).Msg("IPQuality check failed")
			return models.IPQualityResult{IP: proxyStr, Status: "Dead", Error: err.Error()}
		}
		return *res
	})

	for i, p := range proxyStrings {
		wp.AddJob(checker.Job{ID: i, Type: "ipquality", Data: p})
	}

	for i := 0; i < len(proxyStrings); i++ {
		res := <-wp.Results()
		results = append(results, res.(models.IPQualityResult))
	}
	wp.Stop()

	a.logger.Info().Int("results", len(results)).Msg("IPQuality check completed")
	return results
}
