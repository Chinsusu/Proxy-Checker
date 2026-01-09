package proxy

import (
	"context"
	"crypto/tls"
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
	// Normalize proxy string.
	finalProxyStr := strings.TrimSpace(proxyStr)
	scheme := "http"
	if strings.Contains(finalProxyStr, "://") {
		parts := strings.SplitN(finalProxyStr, "://", 2)
		scheme = parts[0]
		finalProxyStr = parts[1]
	}

	// Advanced Parsing for various formats:
	// 1. IP:Port:User:Pass
	// 2. User:Pass@IP:Port
	// 3. IP:Port
	var host string
	var user *url.Userinfo

	if strings.Contains(finalProxyStr, "@") {
		// Format: User:Pass@IP:Port
		parts := strings.SplitN(finalProxyStr, "@", 2)
		userInfo := parts[0]
		host = parts[1]
		if uParts := strings.SplitN(userInfo, ":", 2); len(uParts) == 2 {
			user = url.UserPassword(uParts[0], uParts[1])
		} else {
			user = url.User(userInfo)
		}
	} else {
		parts := strings.Split(finalProxyStr, ":")
		if len(parts) >= 4 {
			// Format: IP:Port:User:Pass
			host = net.JoinHostPort(parts[0], parts[1])
			user = url.UserPassword(parts[2], parts[3])
		} else if len(parts) >= 2 {
			// Format: IP:Port
			host = net.JoinHostPort(parts[0], parts[1])
		} else {
			host = finalProxyStr
		}
	}

	proxyURL := &url.URL{
		Scheme: scheme,
		Host:   host,
		User:   user,
	}

	// Logging (SAFE: Do not modify the original proxyURL object)
	displayUser := "none"
	if proxyURL.User != nil {
		displayUser = proxyURL.User.Username()
	}
	log.Debug().
		Str("scheme", proxyURL.Scheme).
		Str("host", proxyURL.Host).
		Str("user", displayUser).
		Msg("Proxy client initialized (Actual password is kept secret in logs)")

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

func (pc *ProxyClient) RawTCPCheck() error {
	dialer := &net.Dialer{Timeout: pc.Timeout}
	conn, err := dialer.Dial("tcp", pc.Proxy.Host)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
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
