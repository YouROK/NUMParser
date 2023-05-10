package tmdb

import (
	"NUMParser/db/db"
	"NUMParser/db/models"
	"NUMParser/db/utils"
	"encoding/json"
	"errors"
	bolt "go.etcd.io/bbolt"
	"log"
	"time"
)

func GetMovie(id int64) *models.Entity {
	var ent *models.Entity
	db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("TMDB"))
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte("Ents"))
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte("Movies"))
		if bucket == nil {
			return nil
		}
		buf := bucket.Get(utils.I2B(id))
		if len(buf) == 0 {
			return nil
		}
		return json.Unmarshal(buf, &ent)
	})
	return ent
}

func GetTV(id int64) *models.Entity {
	var ent *models.Entity
	db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("TMDB"))
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte("Ents"))
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte("TV"))
		if bucket == nil {
			return nil
		}
		buf := bucket.Get(utils.I2B(id))
		if len(buf) == 0 {
			return nil
		}
		return json.Unmarshal(buf, &ent)
	})
	return ent
}

func FindIMDB(id string) *models.Entity {
	var ent *models.Entity
	db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("TMDB"))
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte("Ents"))
		if bucket == nil {
			return nil
		}
		bckt := bucket.Bucket([]byte("TV"))
		if bckt == nil {
			return nil
		}
		bckt.ForEach(func(_, v []byte) error {
			var e *models.Entity
			err := json.Unmarshal(v, &e)
			if err != nil {
				log.Fatalln("Error read from db TMDB:", err)
			}
			if e.ImdbID == id {
				ent = e
				return errors.New("")
			}
			return nil
		})
		bckt = bucket.Bucket([]byte("Movies"))
		if bckt == nil {
			return nil
		}
		bckt.ForEach(func(_, v []byte) error {
			var e *models.Entity
			err := json.Unmarshal(v, &e)
			if err != nil {
				log.Fatalln("Error read from db TMDB:", err)
			}
			if e.ImdbID == id {
				ent = e
				return errors.New("")
			}
			return nil
		})
		return nil
	})
	return ent
}

func AddTMDB(t *models.Entity) {
	t.UpdateDate = time.Now()
	db.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("TMDB"))
		if err != nil {
			return err
		}
		bucket, err = tx.CreateBucketIfNotExists([]byte("Ents"))
		if err != nil {
			return err
		}
		if t.MediaType == "movie" {
			bucket, err = tx.CreateBucketIfNotExists([]byte("Movies"))
		} else {
			bucket, err = tx.CreateBucketIfNotExists([]byte("TV"))
		}
		if err != nil {
			return err
		}
		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}
		return bucket.Put(utils.I2B(t.ID), buf)
	})
}
