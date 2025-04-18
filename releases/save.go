package releases

import (
	"NUMParser/config"
	"NUMParser/db/models"
	"NUMParser/utils"
	"compress/gzip"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func save(fname string, ents []*models.Entity) {
	if len(ents) == 0 {
		return
	}
	ents = utils.Filter(ents, func(i int, e *models.Entity) bool {
		return e == nil || e.GetTorrent() == nil
	})
	ents = utils.Distinct(ents, func(e *models.Entity) string {
		return strconv.FormatInt(e.ID, 10)
	})
	rid := new(ReleasesID)
	rid.Date = time.Now().Format("02.01.2006")
	rid.Time = time.Now().Format("15:04:05")
	for _, e := range ents {
		if e != nil && e.GetTorrent() != nil {
			var countries []string
			if len(e.ProductionCountries) > 0 {
				for _, c := range e.ProductionCountries {
					countries = append(countries, c.Iso31661)
				}
			} else {
				countries = e.OriginCountry
			}
			d := e.GetTorrent()
			t := Torrent{
				Name:     d.Name,
				Date:     d.CreateDate.Format("02.01.2006"),
				Magnet:   d.Magnet,
				Size:     d.Size,
				Upload:   strconv.Itoa(d.Seed),
				Download: strconv.Itoa(d.Peer),
				Source:   "Rutor",
				Link:     d.Link,
				Quality:  d.VideoQuality,
				Voice:    d.AudioQuality,
			}
			tid := &TmdbId{
				Id:          e.ID,
				MediaType:   e.MediaType,
				GenreIds:    e.GenresIds,
				VoteAverage: e.VoteAverage,
				VoteCount:   e.VoteCount,
				Countries:   countries,
				ReleaseDate: e.ReleaseDate,
				Torrent:     []*Torrent{&t},
			}
			rid.Items = append(rid.Items, tid)
		} else {
			if e == nil {
				log.Println("Error save, empty ent")
			} else if e.GetTorrent() == nil {
				log.Println("Error save, empty torrent:", e.Title)
			}
		}
	}

	if len(rid.Items) > 0 {
		os.MkdirAll(config.SaveReleasePath, 0777)
		err := zip(rid, filepath.Join(config.SaveReleasePath, fname))
		if err != nil {
			log.Println("Error save:", err)
		}
	} else {
		log.Println("**********************************************************")
		log.Println("Empty rid", rid)
		log.Println("Len ents", len(ents))
		log.Println("**********************************************************")
	}
}

func zip(rid *ReleasesID, fname string) error {
	ff, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer ff.Close()
	zw := gzip.NewWriter(ff)
	defer zw.Close()

	enc := json.NewEncoder(zw)
	err = enc.Encode(rid)

	return err
}
