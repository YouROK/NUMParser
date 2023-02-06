package db

import (
	"NUMParser/db/models"
	"NUMParser/db/torrsearch"
	"bytes"
	"compress/flate"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	torrs         []*models.TorrentDetails
	IsTorrsChange bool
	muTorrs       sync.Mutex

	tmdbs         []*models.Entity
	isTMDBsChange bool
	muTmdbs       sync.Mutex

	indxs   map[string]int64
	muIndxs sync.Mutex

	kpds         []*models.KPDetail
	isKPDSChange bool
	muKpds       sync.Mutex
)

func Init() {
	log.Println("Read cache...")
	dir := filepath.Dir(os.Args[0])
	buf, err := os.ReadFile(filepath.Join(dir, "rutor.ls"))
	if err == nil {
		r := flate.NewReader(bytes.NewReader(buf))
		buf, err = io.ReadAll(r)
		r.Close()
		if err == nil {
			var ftors []*models.TorrentDetails
			err = json.Unmarshal(buf, &ftors)
			if err == nil {
				torrs = ftors
				torrsearch.NewIndex(GetTorrs())
			}
		}
	}

	buf, err = os.ReadFile(filepath.Join(dir, "tmdbs.ls"))
	if err == nil {
		r := flate.NewReader(bytes.NewReader(buf))
		buf, err = io.ReadAll(r)
		r.Close()
		if err == nil {
			var ents []*models.Entity
			err = json.Unmarshal(buf, &ents)
			if err == nil {
				tmdbs = ents
			}
		}
	}

	buf, err = os.ReadFile(filepath.Join(dir, "indxs.ls"))
	if err == nil {
		r := flate.NewReader(bytes.NewReader(buf))
		buf, err = io.ReadAll(r)
		r.Close()
		if err == nil {
			var ind map[string]int64
			err = json.Unmarshal(buf, &ind)
			if err == nil {
				indxs = ind
			}
		}
	}
	if indxs == nil {
		indxs = map[string]int64{}
	}

	//buf, err = os.ReadFile(filepath.Join(dir,"kpds.ls"))
	//if err == nil {
	//	r := flate.NewReader(bytes.NewReader(buf))
	//	buf, err = io.ReadAll(r)
	//	r.Close()
	//	if err == nil {
	//		var kpds []*models.KPDetail
	//		err = json.Unmarshal(buf, &kpds)
	//		if err == nil {
	//			SetKPDetails(kpds)
	//		}
	//	}
	//}

	//go savePeriodic()
}

func savePeriodic() {
	for {
		time.Sleep(time.Second * 60 * 10)
		SaveAll()
	}
}

func SaveAll() {
	SaveTorrs()
	SaveTMDB()
	SaveIndxs()
}
