package releases

import (
	"NUMParser/config"
	"NUMParser/db"
	"NUMParser/db/models"
	"NUMParser/utils"
	"log"
	"sort"
	"time"
)

func GetNewCartoons() {
	torrs := db.GetTorrs()
	var list []*models.TorrentDetails

	for _, torr := range torrs {
		if torr.Categories == models.CatCartoonMovie && utils.Abs(torr.Year-time.Now().Year()) < 2 {
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

	ents := FillTMDB("Cartoons", true, list)

	log.Println("Found torrents:", len(ents))
	log.Println("All torrents:", len(list))

	save("cartoons_id.json", ents)
	utils.FreeOSMemGC()
}
