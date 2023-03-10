package tmdb

import (
	"NUMParser/db"
	"NUMParser/db/models"
	"NUMParser/utils"
	"github.com/jmcvetta/napping"
	"log"
	"sort"
	"strconv"
	"sync"
)

const (
	tmdbEndpoint  = "https://api.themoviedb.org/3/"
	imageEndpoint = "http://image.tmdb.org/t/p/"
)

var (
	apiKey = "45ddf563ac3fb845f2d5c363190d1a33"
	apiMu  sync.Mutex

	genres []*models.Genre
)

func Init() {
	log.Println("Init tmdb")
	lstmg := GetGenres("movie")
	lsttvg := GetGenres("tv")

	if lstmg == nil && lsttvg == nil {
		return
	}

	genres = append(lstmg, lsttvg...)

	sort.Slice(genres, func(i, j int) bool {
		return genres[i].Name < genres[j].Name
	})

}

func GetVideoDetails(isMovie bool, id int64) *models.Entity {

	ents := db.GetTMDBDetails()
	for _, ent := range ents {
		if isMovie && ent.MediaType == "movie" && ent.ID == id {
			return ent
		}
		if !isMovie && ent.MediaType == "tv" && ent.ID == id {
			return ent
		}
	}

	params := map[string]string{}
	params["api_key"] = apiKey

	if _, ok := params["language"]; !ok {
		params["language"] = "ru"
	}

	ids := strconv.FormatInt(id, 10)

	endpoint := ""
	if isMovie {
		endpoint = "movie/" + ids
	} else {
		endpoint = "tv/" + ids
	}

	pageParams := napping.Params{}
	for k, v := range params {
		pageParams[k] = v
	}

	var ent *models.Entity

	err := readPageTmdb(endpoint, params, &ent)
	if err != nil || ent == nil {
		return nil
	}
	fixEntity(ent)

	titles := alternativeTitles(isMovie, id)
	ent.Titles = titles

	db.AddTMDB(ent)

	return ent
}

func Search(isMovie bool, query string) []*models.Entity {
	var st = "movie"
	if !isMovie {
		st = "tv"
	}

	params := map[string]string{}
	params["query"] = query

	return listVideoPages("search/"+st, params)
}

func FindByID(isMovie bool, id string, idType string) *models.Entity {

	list := db.GetTMDBDetails()
	for _, e := range list {
		if e.ImdbID == id {
			return e
		}
	}

	params := napping.Params{}

	params["api_key"] = apiKey
	params["external_source"] = idType
	params["language"] = "ru"

	var results *models.FindResult

	err := readPageTmdb("find/"+id, params, &results)
	if err != nil {
		return nil
	}

	if results == nil {
		return nil
	}

	var ent *models.Entity
	if isMovie {
		if len(results.MovieResults) > 0 {
			ent = results.MovieResults[0]
		}
	} else {
		if len(results.TVResults) > 0 {
			ent = results.TVResults[0]
		}
	}

	if ent == nil {
		return nil
	}

	ent = GetVideoDetails(isMovie, ent.ID)
	if ent == nil {
		return nil
	}
	return ent
}

func Legends() []*models.Entity {
	list := listVideoPages("movie/top_rated", map[string]string{})
	for _, e := range list {
		names := alternativeTitles(true, e.ID)
		e.Titles = names
	}
	return list
}

func alternativeTitles(isMovie bool, id int64) []string {
	params := napping.Params{}
	params["api_key"] = apiKey

	var st = "movie"
	if !isMovie {
		st = "tv"
	}
	var results *models.AlternativeTitles

	err := readPageTmdb(st+"/"+strconv.FormatInt(id, 10)+"/alternative_titles", params, &results)
	if err != nil {
		return nil
	}

	if results == nil {
		return nil
	}

	var list []string

	for _, title := range results.Titles {
		if title.Title != "" {
			list = append(list, title.Title)
		}
	}

	return list
}

func listVideoPages(endpoint string, params napping.Params) []*models.Entity {
	p := map[string]string{}
	for k, v := range params {
		p[k] = v
	}
	p["page"] = "1"
	lst, pages := listVideo(endpoint, p)
	if pages > 10 {
		pages = 10
	}
	if pages > 1 {
		lsts := make([][]*models.Entity, pages-1)
		utils.ParallelFor(2, pages+1, func(i int) {
			p := map[string]string{}
			for k, v := range params {
				p[k] = v
			}
			p["page"] = strconv.Itoa(i)
			lsts[i-2], _ = listVideo(endpoint, p)
		})
		for _, l := range lsts {
			lst = append(lst, l...)
		}
	}

	return lst
}

func listVideo(endpoint string, params napping.Params) ([]*models.Entity, int) {
	params["api_key"] = apiKey

	if _, ok := params["language"]; !ok {
		params["language"] = "ru"
	}

	var results *models.EntityRequest
	pageParams := napping.Params{}
	for k, v := range params {
		pageParams[k] = v
	}

	err := readPageTmdb(endpoint, params, &results)
	if err != nil {
		return nil, 0
	}

	if results == nil {
		return nil, 0
	}

	for _, v := range results.Results {
		fixEntity(v)
	}

	return results.Results, results.TotalPages
}
