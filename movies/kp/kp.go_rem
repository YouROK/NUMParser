package kp

import (
	"NUMParser/db"
	"NUMParser/db/models"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	apiKey = ""
	host   = "https://kinopoiskapiunofficial.tech"
	kpList = []*models.KPDetail{}
)

func GetDetail(id string) (*models.KPDetail, error) {
	var body []byte
	st := 0
	var err error
	for i := 0; i < 20; i++ {
		body, st, err = get(host + "/api/v2.2/films/" + id)
		if st == 429 {
			time.Sleep(time.Second)
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}

	var kpd *models.KPDetail
	err = json.Unmarshal(body, &kpd)
	if err == nil {
		if kpd.KinopoiskID == 1103803 {
			kpd.Year = 2022
		}
		db.AddKPD(kpd)
	}
	return kpd, err
}

func Search(query string) ([]*models.KPDetail, error) {
	var body []byte
	st := 0
	var err error
	query = url.QueryEscape(query)
	for i := 0; i < 20; i++ {
		body, st, err = get(host + "/api/v2.1/films/search-by-keyword?keyword=" + query + "&page=1")
		if st == 429 {
			time.Sleep(time.Second)
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}

	var kpd *models.KPSearch
	err = json.Unmarshal(body, &kpd)

	var list []*models.KPDetail
	dblist := db.GetKPDetails()
	for _, film := range kpd.Films {
		isFind := false
		for _, d := range dblist {
			if d.KinopoiskID == film.FilmID {
				list = append(list, d)
				isFind = true
				break
			}
		}
		if !isFind {
			kid := strconv.Itoa(film.FilmID)
			d, err := GetDetail(kid)
			if err == nil {
				list = append(list, d)
			}
		}
	}
	return list, err
}

func get(link string) ([]byte, int, error) {
	var httpClient *http.Client
	httpClient = &http.Client{
		Timeout: 120 * time.Second,
	}
	req, err := http.NewRequest("GET", link, nil)

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	req.Header.Set("X-API-KEY", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("Error get link:", link, resp.StatusCode, resp.Status)
		return nil, resp.StatusCode, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	return body, 0, nil
}
