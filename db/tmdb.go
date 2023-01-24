package db

import (
	"NUMParser/db/models"
	"bytes"
	"compress/flate"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

func GetTMDBDetails() []*models.Entity {
	return tmdbs
}

func SetTMDBDetails(list []*models.Entity) {
	tmdbs = list
	isTMDBsChange = true
}

func AddTMDB(t *models.Entity) {
	isTMDBsChange = true
	t.UpdateDate = time.Now()
	for i, d := range tmdbs {
		if d.ID == t.ID {
			tmdbs[i] = t
			return
		}
	}

	muTmdbs.Lock()
	defer muTmdbs.Unlock()
	tmdbs = append(tmdbs, t)
}

func SaveTMDB() {
	muTmdbs.Lock()
	var list []*models.Entity
	for _, kpd := range tmdbs {
		if time.Now().Before(kpd.UpdateDate.Add(time.Hour * 24 * 7)) {
			list = append(list, kpd)
		}
	}
	if len(list) != len(tmdbs) {
		tmdbs = list
		isTMDBsChange = true
	}
	muTmdbs.Unlock()

	if !isTMDBsChange || len(tmdbs) == 0 {
		return
	}

	log.Println("Save TMDB")
	buf, err := json.Marshal(tmdbs)
	if err != nil {
		return
	}

	var b bytes.Buffer
	w, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}
	w.Write(buf)
	w.Close()

	dir := filepath.Dir(os.Args[0])
	err = os.WriteFile(filepath.Join(dir, "tmdbs.ls"), b.Bytes(), 0666)
	if err != nil {
		return
	}
	isTMDBsChange = false
}
