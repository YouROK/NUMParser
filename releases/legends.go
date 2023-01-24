package releases

import (
	"NUMParser/db"
	"NUMParser/db/models"
	"NUMParser/movies/tmdb"
	"NUMParser/utils"
	"log"
	"sort"
)

func GetLegends() {
	log.Println("Search legends")
	ents := tmdb.Legends()
	ents = filterCyrilic(ents)
	if len(ents) > 520 {
		ents = ents[:520]
	}
	var found []*models.Entity
	for i, e := range ents {
		list := db.SearchTorr(e.Title + " " + e.OriginalTitle + " " + e.Year)
		if len(list) == 0 {
			names := getEnNames(e.Titles)
			for _, name := range names {
				list = db.SearchTorr(e.Title + " " + name + " " + e.Year)
				if len(list) > 0 {
					break
				}
			}
		}
		if len(list) > 0 {
			sort.Slice(list, func(i, j int) bool {
				if list[i].CreateDate == list[j].CreateDate {
					if list[i].VideoQuality == list[j].VideoQuality {
						return list[i].AudioQuality > list[j].AudioQuality
					}
					return list[i].VideoQuality > list[j].VideoQuality
				}
				return list[i].CreateDate.After(list[j].CreateDate)
			})

			e.SetTorrent(list[0])
			found = append(found, e)
		}
		log.Println("Fill legends:", i+1, "/", len(ents))
		utils.FreeOSMemGC()
	}

	log.Println("Found torrents:", len(found))
	log.Println("All torrents:", len(ents))

	save("legends_id.json", found)

	utils.FreeOSMemGC()
}

func getEnNames(names []string) []string {
	var list []string
	for i, name := range names {
		if len(utils.ClearStr(name)) > 0 {
			list = append(list, names[i])
		}
	}
	return list
}

func filterCyrilic(ents []*models.Entity) []*models.Entity {
	return utils.Filter(ents, func(i int, e *models.Entity) bool {
		str := utils.ClearStr(e.Title)
		return len(str) == 0
	})
}
