package db

import (
	"NUMParser/db/models"
	"compress/flate"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetTorrs() []*models.TorrentDetails {
	return torrs
}

func SetTorrs(list []*models.TorrentDetails) {
	torrs = list
	IsTorrsChange = true
}

func GetTorrsCategory(cat string) []*models.TorrentDetails {
	muTorrs.Lock()
	defer muTorrs.Unlock()
	var list []*models.TorrentDetails
	for _, torr := range torrs {
		if torr.Categories == cat {
			list = append(list, torr)
		}
	}
	return list
}

func AddTorr(t *models.TorrentDetails) {
	if t.Hash != "" {
		for i, tdb := range torrs {
			if tdb.Hash == t.Hash {
				t.IMDBID = torrs[i].IMDBID
				torrs[i] = t
				return
			}
		}
	}

	muTorrs.Lock()
	defer muTorrs.Unlock()
	IsTorrsChange = true
	torrs = append(torrs, t)
}

func SaveTorrs() {
	removeOldTorr()
	if !IsTorrsChange || len(torrs) == 0 {
		return
	}
	log.Println("Save torrents")

	dir := filepath.Dir(os.Args[0])
	ff, err := os.Create(filepath.Join(dir, "rutor.ls"))
	if err != nil {
		log.Println("Error save torrs:", err)
		return
	}
	defer ff.Close()

	w, err := flate.NewWriter(ff, flate.BestCompression)
	if err != nil {
		log.Println("Error save torrs:", err)
		return
	}
	defer w.Close()

	enc := json.NewEncoder(w)
	err = enc.Encode(torrs)
	if err != nil {
		log.Println("Error save torrs:", err)
		return
	}
	indexTorrs()
	IsTorrsChange = false
}

func RemoveSeriesTorr() {
	muTorrs.Lock()
	defer muTorrs.Unlock()

	var list []*models.TorrentDetails

	for _, t := range torrs {
		if t.Categories != models.CatSeries {
			list = append(list, t)
		}
	}

	torrs = list
}

func removeOldTorr() {
	muTorrs.Lock()
	defer muTorrs.Unlock()

	var list []*models.TorrentDetails
	inResult := make(map[string]bool)

	for _, t := range torrs {
		link := t.Link
		if _, ok := inResult[link]; !ok {
			inResult[link] = true
			list = append(list, t)
		}
	}

	torrs = list
}
