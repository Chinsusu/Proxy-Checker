package proxy

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/proxy"
)

type ProxyClient struct {
	Proxy      *url.URL
	HTTPClient *http.Client
	UserAgent  string
	Timeout    time.Duration
}

func NewProxyClient(proxyStr string, userAgent string, timeout time.Duration) (*ProxyClient, error) {
	// Normalize proxy string. Default to http if no scheme is provided.
	finalProxyStr := proxyStr
	scheme := "http"
	if strings.Contains(finalProxyStr, "://") {
		parts := strings.SplitN(finalProxyStr, "://", 2)
		scheme = parts[0]
		finalProxyStr = parts[1]
	}

	// Handle IP:Port:User:Pass format
	parts := strings.Split(finalProxyStr, ":")
	var proxyURL *url.URL
	var err error

	if len(parts) >= 4 {
		// IP:Port:User:Pass -> scheme://User:Pass@IP:Port
		u := fmt.Sprintf("%s://%s:%s@%s:%s", scheme, parts[2], parts[3], parts[0], parts[1])
		proxyURL, err = url.Parse(u)
	} else {
		// IP:Port or scheme://IP:Port
		u := fmt.Sprintf("%s://%s", scheme, finalProxyStr)
		proxyURL, err = url.Parse(u)
	}

	if err != nil {
		return nil, err
	}

	// Logging (sanitized)
	sanitizedURL := *proxyURL
	if sanitizedURL.User != nil {
		sanitizedURL.User = url.UserPassword(sanitizedURL.User.Username(), "********")
	}
	log.Debug().Str("proxy_parsed", sanitizedURL.String()).Msg("NewProxyClient created")

	baseDialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// proxy.FromURL only supports socks5, socks4, socks4a.
	// For http/https, we use the standard transport.Proxy.
	if strings.HasPrefix(proxyURL.Scheme, "http") {
		transport.Proxy = http.ProxyURL(proxyURL)
		transport.DialContext = baseDialer.DialContext
	} else {
		dialer, err := proxy.FromURL(proxyURL, baseDialer)
		if err != nil {
			return nil, err
		}
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return &ProxyClient{
		Proxy:      proxyURL,
		HTTPClient: client,
		UserAgent:  userAgent,
		Timeout:    timeout,
	}, nil
}

func (pc *ProxyClient) TestConnectivity() (bool, error) {
	resp, err := pc.HTTPClient.Get("https://www.google.com")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

func GetRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	}
	return userAgents[time.Now().UnixNano()%int64(len(userAgents))]
}
