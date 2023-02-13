package releases

import (
	"NUMParser/config"
	"NUMParser/db"
	"NUMParser/db/models"
	"NUMParser/utils"
	"log"
	"math"
	"sort"
	"strconv"
	"time"
)

func GetNewMovies() {
	torrs := db.GetTorrs()
	var list []*models.TorrentDetails

	for _, torr := range torrs {
		if torr.Categories == models.CatMovie && utils.Abs(torr.Year-time.Now().Year()) < 2 && torr.VideoQuality >= 200 {
			list = append(list, torr)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].CreateDate == list[j].CreateDate {
			if list[i].VideoQuality == list[j].VideoQuality {
				return list[i].AudioQuality > list[j].AudioQuality
			}
			return list[i].VideoQuality > list[j].VideoQuality
		}
		return list[i].CreateDate.After(list[j].CreateDate)
	})

	list = utils.UniqueTorrList(list)

	if config.ReleasesLimit > 0 && len(list) > config.ReleasesLimit {
		list = list[:config.ReleasesLimit]
	}

	ents := FillTMDB("Movies", true, list)

	log.Println("Found torrents:", len(ents))
	log.Println("All torrents:", len(list))

	save("movies_id.json", ents)
	utils.FreeOSMemGC()
}

func GetNewMoviesYear(year int) {
	torrs := db.GetTorrs()
	var list []*models.TorrentDetails

	for _, torr := range torrs {
		if (torr.Categories == models.CatMovie || torr.Categories == models.CatCartoonMovie) && torr.Year == year && torr.VideoQuality >= 100 {
			list = append(list, torr)
		}
	}

	//sort.Slice(list, func(i, j int) bool {
	//	if list[i].CreateDate == list[j].CreateDate {
	//		return list[i].VideoQuality > list[j].VideoQuality
	//	}
	//	return list[i].CreateDate.After(list[j].CreateDate)
	//})

	list = utils.UniqueTorrList(list)

	if config.ReleasesLimit > 0 && len(list) > config.ReleasesLimit {
		list = list[:config.ReleasesLimit]
	}

	ents := FillTMDB("Movies "+strconv.Itoa(year), true, list)

	ents = utils.Filter(ents, func(i int, e *models.Entity) bool {
		return e == nil || e.GetTorrent() == nil
	})
	ents = utils.Distinct(ents, func(e *models.Entity) string {
		return strconv.FormatInt(e.ID, 10)
	})

	sort.Slice(ents, func(i, j int) bool {
		rankI := ents[i].VoteAverage * math.Log(float64(ents[i].VoteCount))
		rankJ := ents[j].VoteAverage * math.Log(float64(ents[j].VoteCount))

		if math.Abs(rankI-rankJ) < 0.1 {
			return ents[i].Popularity > ents[j].Popularity
		}

		return rankI > rankJ
	})

	log.Println("Found torrents:", len(ents))
	log.Println("All torrents:", len(list))

	save("movies_id_"+strconv.Itoa(year)+".json", ents)
	utils.FreeOSMemGC()
}
