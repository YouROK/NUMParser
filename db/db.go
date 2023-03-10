package db

import (
	"NUMParser/db/models"
	"NUMParser/db/torrsearch"
	"compress/flate"
	"encoding/json"
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
	ff, err := os.Open(filepath.Join(dir, "rutor.ls"))
	if err == nil {
		defer ff.Close()
		r := flate.NewReader(ff)
		defer r.Close()
		if err == nil {
			var ftors []*models.TorrentDetails
			err = json.NewDecoder(r).Decode(&ftors)
			if err == nil {
				torrs = ftors
				torrsearch.NewIndex(GetTorrs())
			}
		}
	}

	ff, err = os.Open(filepath.Join(dir, "tmdbs.ls"))
	if err == nil {
		defer ff.Close()
		r := flate.NewReader(ff)
		defer r.Close()
		if err == nil {
			var ents []*models.Entity
			err = json.NewDecoder(r).Decode(&ents)
			if err == nil {
				tmdbs = ents
			}
		}
	}

	ff, err = os.Open(filepath.Join(dir, "indxs.ls"))
	if err == nil {
		defer ff.Close()
		r := flate.NewReader(ff)
		defer r.Close()
		if err == nil {
			var ind map[string]int64
			err = json.NewDecoder(r).Decode(&ind)
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
