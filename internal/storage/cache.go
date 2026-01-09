package storage

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Cache struct {
	db *sql.DB
}

func NewCache(dbPath string) (*Cache, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value TEXT,
		expires_at DATETIME
	)`)
	if err != nil {
		return nil, err
	}

	return &Cache{db: db}, nil
}

func (c *Cache) Get(key string) (string, bool) {
	var value string
	var expiresAt time.Time
	err := c.db.QueryRow("SELECT value, expires_at FROM cache WHERE key = ?", key).Scan(&value, &expiresAt)
	if err != nil {
		return "", false
	}

	if time.Now().After(expiresAt) {
		c.db.Exec("DELETE FROM cache WHERE key = ?", key)
		return "", false
	}

	return value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	valBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(ttl)
	_, err = c.db.Exec("INSERT OR REPLACE INTO cache (key, value, expires_at) VALUES (?, ?, ?)", key, string(valBytes), expiresAt)
	return err
}

func (c *Cache) Close() error {
	return c.db.Close()
}
