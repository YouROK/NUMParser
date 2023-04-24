package client

import (
	"NUMParser/config"
	"context"
	"errors"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func GetNic(link, referer, cookie string) (string, error) {
	var (
		dnsResolverIP        = "151.80.222.79:53"
		dnsResolverProto     = "udp"
		dnsResolverTimeoutMs = 10000
	)

	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Duration(dnsResolverTimeoutMs) * time.Millisecond,
				}
				return d.DialContext(ctx, dnsResolverProto, dnsResolverIP)
			},
		},
		Timeout:   120 * time.Second,
		KeepAlive: 120 * time.Second,
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", link, nil)

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	if cookie != "" {
		req.Header.Set("cookie", cookie)
	}
	if referer != "" {
		req.Header.Set("referer", referer)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Get(link string) (string, error) {
	buf, err := GetBuf(link, "", "")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func GetTor(link string) (string, error) {
	proxyUrl := "192.168.1.2:9050"
	dialer, err := proxy.SOCKS5("tcp", proxyUrl, nil, proxy.Direct)
	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}
	transport := &http.Transport{DialContext: dialContext,
		DisableKeepAlives: true}
	cl := &http.Client{Transport: transport}

	resp, err := cl.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetBuf(link, referer, cookie string) ([]byte, error) {
	var httpClient *http.Client
	req, err := http.NewRequest("GET", link, nil)

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	if cookie != "" {
		req.Header.Set("cookie", cookie)
	}
	if referer != "" {
		req.Header.Set("referer", referer)
	}

	if config.ProxyHost != "" {
		proxyURL, _ := url.Parse(config.ProxyHost)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		httpClient = &http.Client{
			Transport: transport,
			Timeout:   120 * time.Second,
		}
	} else {
		httpClient = &http.Client{
			Timeout: 120 * time.Second,
		}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Error get link:", link, resp.StatusCode, resp.Status)
		return nil, errors.New(resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
