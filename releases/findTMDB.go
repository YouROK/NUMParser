package releases

import (
	"NUMParser/db/models"
	tmdb2 "NUMParser/db/tmdb"
	"NUMParser/movies/tmdb"
	"NUMParser/parser"
	"NUMParser/utils"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"sync"
)

func FillTMDB(label string, isMovie bool, torrs []*models.TorrentDetails) []*models.Entity {
	list := make([]*models.Entity, len(torrs))
	var mu sync.Mutex
	utils.PForLim(torrs, 20, func(i int, t *models.TorrentDetails) {
		//for i, t := range torrs {
		var md *models.Entity
		indx := tmdb2.GetIndex(t.Hash)
		if indx != 0 {
			md = tmdb.GetVideoDetails(isMovie, indx)
			if md != nil {
				mu.Lock()
				list[i] = md
				mu.Unlock()
			}
		}
		if md == nil {
			md = FindTMDBID(isMovie, t)
			if md != nil {
				mu.Lock()
				list[i] = md
				mu.Unlock()
			} else {
				md = FindTMDB(isMovie, t)
				if md != nil {
					mu.Lock()
					list[i] = md
					mu.Unlock()
				}
			}
		}
		if md == nil {
			log.Println(label+":", "Torr", i, "/", len(torrs), "not found in TMDB:", t.Title, t.Link)
		} else {
			tmdb2.SetIndex(t, md)
			md.SetTorrent(t)
			log.Println(label+":", "Find torr", i, "/", len(torrs), "in TMDB:", t.Title)
		}
		//}
	})

	return list
}

func FindTMDBID(isMovie bool, torr *models.TorrentDetails) *models.Entity {
	if torr.IMDBID != "" {
		return tmdb.FindByID(isMovie, torr.IMDBID, "imdb_id")
	}
	body := parser.GetBodyLink(torr)
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(body))
	if err != nil {
		return nil
	}

	ids := ""

	doc.Find("table#details").Find("a").Each(func(i int, selection *goquery.Selection) {
		if link, ok := selection.Attr("href"); ok {
			if strings.Contains(link, "www.imdb.com") {
				link = strings.TrimRight(link, "/")
				arr := strings.Split(link, "/")
				if len(arr) > 0 {
					ids = arr[len(arr)-1]
					return
				}
			}
		}
	})
	if ids == "" {
		return nil
	}

	torr.IMDBID = ids

	return tmdb.FindByID(isMovie, ids, "imdb_id")
}

func FindTMDB(isMovie bool, torr *models.TorrentDetails) *models.Entity {
	list := tmdb.Search(isMovie, torr.Name)

	list = utils.Filter(list, func(i int, e *models.Entity) bool {
		if len(e.ReleaseDate) > 6 {
			year, _ := strconv.Atoi(e.ReleaseDate[6:])
			return utils.Abs(year-torr.Year) > 1
		}
		return true
	})

	if len(list) == 1 {
		return list[0]
	}

	for _, name := range torr.Names {
		lst := tmdb.Search(true, name)
		list = append(list, lst...)
	}

	list = utils.Filter(list, func(i int, e *models.Entity) bool {
		if len(e.ReleaseDate) > 6 {
			year, _ := strconv.Atoi(e.ReleaseDate[6:])
			return utils.Abs(year-torr.Year) > 1
		}
		return true
	})

	list = utils.Distinct(list, func(e *models.Entity) string {
		return strconv.FormatInt(e.ID, 10)
	})

	if len(list) == 1 {
		return list[0]
	}

	if len(list) > 0 {
		list = utils.Filter(list, func(i int, e *models.Entity) bool {
			if len(e.ReleaseDate) > 6 {
				year, _ := strconv.Atoi(e.ReleaseDate[6:])
				return utils.Abs(year-torr.Year) > 1
			}
			return true
		})
		names := append([]string{torr.Name}, torr.Names...)
		list = utils.Filter(list, func(i int, e *models.Entity) bool {
			finds := 0
			for _, name := range names {
				if utils.ClearStr(e.Title) == utils.ClearStr(e.OriginalTitle) &&
					utils.ClearStr(name) == utils.ClearStr(e.Title) && len(names) == 1 {
					return false
				} else if utils.ClearStr(e.OriginalTitle) == utils.ClearStr(name) ||
					utils.ClearStr(e.Title) == utils.ClearStr(name) {
					finds++
				}
				if finds > 1 {
					return false
				}
				if utils.ClearStr(e.Title) != utils.ClearStr(e.OriginalTitle) && len(names) > 1 {
					for _, title := range e.Titles {
						if utils.ClearStr(title) == utils.ClearStr(name) {
							return false
						}
					}
				}
			}

			return true
		})
		if len(list) == 1 {
			return list[0]
		}
	}
	return nil
}
