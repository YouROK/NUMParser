package client

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	proxyNum = 0
)

func getProxyFromList() string {
	dir := filepath.Dir(os.Args[0])
	fileName := filepath.Join(dir, "proxy.list")
	buf, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Error load proxy list:", err)
		return ""
	}
	list := strings.Split(string(buf), "\n")
	if proxyNum > len(list) {
		proxyNum = 0
		if len(list) > 0 {
			return list[0]
		}
		return ""
	}
	proxyHost := ""
	if len(list) > 0 {
		proxyHost = strings.TrimSpace(list[proxyNum])
		if !strings.HasPrefix(proxyHost, "http") && !strings.HasPrefix(proxyHost, "socks") {
			proxyHost = "//" + proxyHost
		}
	}
	proxyNum++
	return proxyHost
}
