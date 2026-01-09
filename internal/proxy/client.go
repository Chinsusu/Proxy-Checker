package proxy

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type ProxyClient struct {
	Proxy      *url.URL
	HTTPClient *http.Client
	UserAgent  string
	Timeout    time.Duration
}

func NewProxyClient(proxyStr string, userAgent string, timeout time.Duration) (*ProxyClient, error) {
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
