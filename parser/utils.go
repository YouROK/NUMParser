package parser

import (
	"NUMParser/client"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"log"
	"strings"
	"time"
)

func getHash(magnet string) string {
	//magnet:?xt=urn:btih:1debb44e9e9ac785aaa4c26507534e1357672a22&dn=rutor.info&tr=udp://opentor.net:6969&tr=http://retracker.local/announce
	pos := strings.Index(magnet, "btih:")
	if pos == -1 {
		return ""
	}
	magnet = magnet[pos+5:]
	pos = strings.Index(magnet, "&")
	if pos == -1 {
		return strings.ToLower(magnet)
	}
	return strings.ToLower(magnet[:pos])
}

func get(link string) (string, error) {
	var body string
	var err error
	for i := 0; i < 10; i++ {
		if strings.Contains(link, "rutor.lib") {
			body, err = client.GetNic(link, "", "")
		} else {
			body, err = client.Get(link)
		}
		if err == nil {
			break
		}
		log.Println("Error get page,tryes:", i+1, link, err)
		if i < 5 {
			time.Sleep(time.Minute)
		} else {
			time.Sleep(time.Minute * 10)
		}
	}
	return body, err
}

func getBuf(link, referer string) ([]byte, error) {
	var body []byte
	var err error
	for i := 0; i < 10; i++ {
		body, err = client.GetBuf(link, referer, "")
		if err == nil {
			break
		}
		log.Println("Error get page,tryes:", i+1, link)
		time.Sleep(time.Second * 2)
	}
	return body, err
}

func node2Text(node *html.Node) string {
	return strings.TrimSpace(strings.Replace((&goquery.Selection{Nodes: []*html.Node{node}}).Text(), "\u00A0", " ", -1))
}

func replaceBadName(name string) string {
	name = strings.ReplaceAll(name, "Ванда/Вижн ", "ВандаВижн ")
	name = strings.ReplaceAll(name, "Ё", "Е")
	name = strings.ReplaceAll(name, "ё", "е")
	name = strings.ReplaceAll(name, "щ", "ш")
	return name
}
