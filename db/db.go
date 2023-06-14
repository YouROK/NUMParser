package db

import (
	"NUMParser/db/db"
	"NUMParser/db/models"
	"NUMParser/db/rutor"
	"log"
)

func Init() {
	log.Println("Open db...")
	db.Init()

	log.Println("Read db...")
	rutor.Init()
}

func SaveAll() {
	rutor.SaveTorrs()
}

func GetTorrs() []*models.TorrentDetails {
	var torrs []*models.TorrentDetails
	torrs = append(torrs, rutor.GetTorrs()...)
	return torrs
}
