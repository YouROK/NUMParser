package rutor

import (
	"NUMParser/db/db"
	"NUMParser/db/models"
	"NUMParser/db/torrsearch"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"log"
	"strings"
	"sync"
)

var (
	torrs         []*models.TorrentDetails
	IsTorrsChange bool
	muTorrs       sync.Mutex
)

func Init() {
	db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Rutor"))
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte("Torrents"))
		if bucket == nil {
			return nil
		}
		err := bucket.ForEach(func(_, v []byte) error {
			var torr *models.TorrentDetails
			err := json.Unmarshal(v, &torr)
			if err == nil {
				torrs = append(torrs, torr)
			}
			return err
		})
		if err != nil {
			log.Println("Error read rutor from db:", err)
		}
		return nil
	})

	torrsearch.NewIndex(torrs)
}

func RemoveAll() {
	torrs = nil
	db.DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("Rutor"))
		return nil
	})
}

func GetTorrs() []*models.TorrentDetails {
	return torrs
}

func SetTorrs(list []*models.TorrentDetails) {
	torrs = list
	IsTorrsChange = true
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
	muTorrs.Lock()
	defer muTorrs.Unlock()
	log.Println("Save torrents")

	err := db.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("Rutor"))
		if err != nil {
			return err
		}
		//Recreate torrents
		bucket.DeleteBucket([]byte("Torrents"))
		bucket, err = bucket.CreateBucket([]byte("Torrents"))
		if err != nil {
			return err
		}

		for _, torr := range torrs {
			buf, err := json.Marshal(torr)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(strings.ToLower(torr.Hash)), buf)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalln("Error write to db rutor:", err)
	}
	torrsearch.NewIndex(torrs)
	IsTorrsChange = false
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
