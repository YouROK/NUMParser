package db

import (
	"NUMParser/db/models"
	"bytes"
	"compress/flate"
	"encoding/json"
	"log"
	"os"
)

func GetIndex(hash string) int64 {
	muIndxs.Lock()
	defer muIndxs.Unlock()
	if i, ok := indxs[hash]; ok {
		return i
	}

	return 0
}

func SetIndex(t *models.TorrentDetails, e *models.Entity) {
	muIndxs.Lock()
	defer muIndxs.Unlock()
	indxs[t.Hash] = e.ID
}

func GetIndexes() map[string]int64 {
	return indxs
}

func SaveIndxs() {
	muIndxs.Lock()
	defer muIndxs.Unlock()
	log.Println("Save Index")
	buf, err := json.Marshal(indxs)
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

	os.WriteFile("indxs.ls", b.Bytes(), 0666)
}
