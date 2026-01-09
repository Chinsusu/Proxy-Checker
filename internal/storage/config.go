package storage

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	API struct {
		IPWho struct {
			RateLimit int           `yaml:"rate_limit"`
			Timeout   time.Duration `yaml:"timeout"`
			CacheTTL  time.Duration `yaml:"cache_ttl"`
		} `yaml:"ipwho"`
		IPQuality struct {
			Timeout    time.Duration `yaml:"timeout"`
			UserAgents []string      `yaml:"user_agents"`
			DelayRange []int         `yaml:"delay_range"`
		} `yaml:"ipquality"`
	} `yaml:"api"`
	Worker struct {
		PoolSize      int           `yaml:"pool_size"`
		RetryAttempts int           `yaml:"retry_attempts"`
		RetryDelay    time.Duration `yaml:"retry_delay"`
	} `yaml:"worker"`
	Proxy struct {
		ConnectionTimeout time.Duration `yaml:"connection_timeout"`
		Types             []string      `yaml:"types"`
	} `yaml:"proxy"`
	Storage struct {
		CacheEnabled bool   `yaml:"cache_enabled"`
		DBPath       string `yaml:"db_path"`
	} `yaml:"storage"`
}

var DefaultConfig = `
api:
  ipwho:
    rate_limit: 10
    timeout: 10s
    cache_ttl: 24h
  ipquality:
    timeout: 30s
    user_agents:
      - "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
    delay_range: [1000, 3000]

worker:
  pool_size: 5
  retry_attempts: 3
  retry_delay: 2s

proxy:
  connection_timeout: 10s
  types: ["http", "https", "socks5"]

storage:
  cache_enabled: true
  db_path: "./data/cache.db"
`

func LoadConfig(path string) (*Config, error) {
	var config Config
	f, err := os.Open(path)
	if err != nil {
		// If fails to load, use default or return error
		err = yaml.Unmarshal([]byte(DefaultConfig), &config)
		return &config, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
