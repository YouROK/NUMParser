package client

import (
	"NUMParser/config"
	"context"
	"errors"
	"github.com/parnurzeal/gorequest"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var Err404 = errors.New("404 not found")

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
	if resp.StatusCode == 404 {
		log.Println("Error get link:", link, resp.StatusCode, resp.Status)
		return "", Err404
	} else if resp.StatusCode != 200 {
		log.Println("Error get link:", link, resp.StatusCode, resp.Status)
		return "", errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Get(link string) *gorequest.SuperAgent {
	return GetParam(link, "", "")
}

func GetParam(link, referer, cookie string) *gorequest.SuperAgent {
	agent := gorequest.New()

	if cookie != "" {
		header := http.Header{}
		header.Add("Cookie", cookie)
		request := http.Request{
			Header: header,
		}
		agent.Cookies = request.Cookies()
	}

	if referer != "" {
		agent.AppendHeader("referer", referer)
	}

	if config.UseProxy {
		proxyHost := getProxyFromList()
		if proxyHost == "" {
			proxyHost = config.ProxyHost
		}
		if proxyHost != "" {
			agent.Proxy(proxyHost)
		}
	} else if config.ProxyHost != "" {
		agent.Proxy(config.ProxyHost)
	}
	agent.Timeout(30 * time.Second)
	return agent.Get(link)
}
