package db

import (
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
	"path/filepath"
)

var (
	DB *bolt.DB
)

func Init() {
	dir := filepath.Dir(os.Args[0])
	dir = filepath.Join(dir, "db")
	os.MkdirAll(dir, 0777)
	fileDb := filepath.Join(dir, "numparser.db")
	var err error
	DB, err = bolt.Open(fileDb, 0666, nil)
	if err != nil {
		log.Fatalln("Error open db", err)
	}
}
