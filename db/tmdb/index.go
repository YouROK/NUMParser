package tmdb

import (
	"NUMParser/db/db"
	"NUMParser/db/models"
	"NUMParser/db/utils"
	bolt "go.etcd.io/bbolt"
	"log"
)

func GetIndex(hash string) int64 {
	var id int64
	db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("TMDB"))
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte("Index"))
		if bucket == nil {
			return nil
		}
		buf := bucket.Get([]byte(hash))
		if len(buf) == 0 {
			return nil
		}
		id = utils.B2I(buf)
		return nil
	})
	return id
}

func SetIndex(t *models.TorrentDetails, e *models.Entity) {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("TMDB"))
		if err != nil {
			return err
		}
		bucket, err = bucket.CreateBucketIfNotExists([]byte("Index"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(t.Hash), utils.I2B(e.ID))
	})
	if err != nil {
		log.Fatalln("Error write to db TMDB:", err)
	}
}
