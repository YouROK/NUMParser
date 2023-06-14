package tmdb

import (
	"NUMParser/db/models"
	"github.com/jmcvetta/napping"
	"strings"
	"time"
)

func ImageURL(uri string, size string) string {
	if uri == "" {
		return ""
	}
	if strings.HasPrefix(uri, imageEndpoint) {
		return uri
	}
	return imageEndpoint + size + uri
}

func FixDate(date string) string {
	tm, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return tm.Format("02.01.2006")
}

func GetGenres(gtype string) []*models.Genre {
	var genres *models.GenresRequest

	urlPages := napping.Params{
		//"api_key":  apiKey,
		"language": "ru",
	}
	endpoint := "genre/" + gtype + "/list"

	err := readPageTmdb(endpoint, urlPages, &genres)
	if err != nil {
		return nil
	}

	if genres != nil && len(genres.Genres) > 0 {
		for _, i := range genres.Genres {
			i.Name = strings.Title(i.Name)
		}
	} else {
		return nil
	}
	return genres.Genres
}

func fixEntity(ent *models.Entity) {
	if len(ent.ReleaseDate) >= 4 {
		ent.Year = ent.ReleaseDate[:4]
	} else if len(ent.FirstAirDate) >= 4 {
		ent.Year = ent.FirstAirDate[:4]
	}

	if ent.MediaType == "" {
		if ent.Title == "" {
			ent.MediaType = "tv"
		} else if ent.Name == "" {
			ent.MediaType = "movie"
		}
	}

	if ent.Title == "" && ent.Name != "" {
		ent.Title = ent.Name
	}

	if ent.OriginalTitle == "" && ent.OriginalName != "" {
		ent.OriginalTitle = ent.OriginalName
	}

	ent.ReleaseDate = FixDate(ent.ReleaseDate)
	ent.FirstAirDate = FixDate(ent.FirstAirDate)
	ent.LastAirDate = FixDate(ent.LastAirDate)

	if ent.ReleaseDate == "" && ent.FirstAirDate != "" {
		ent.ReleaseDate = ent.FirstAirDate
	}

	ent.PosterPath = ImageURL(ent.PosterPath, "w342")
	ent.BackdropPath = ImageURL(ent.BackdropPath, "w780")

	if ent.BelongsToCollection != nil {
		ent.BelongsToCollection.Poster = ImageURL(ent.BelongsToCollection.Poster, "w342")
		ent.BelongsToCollection.Backdrop = ImageURL(ent.BelongsToCollection.Backdrop, "w780")
	}

	if ent.Images != nil {
		for i := range ent.Images.Backdrops {
			ent.Images.Backdrops[i].FilePath = ImageURL(ent.Images.Backdrops[i].FilePath, "w780")
		}

		for i := range ent.Images.Posters {
			ent.Images.Posters[i].FilePath = ImageURL(ent.Images.Posters[i].FilePath, "w342")
		}
	}

	for i := range ent.ProductionCompanies {
		ent.ProductionCompanies[i].LogoPath = ImageURL(ent.ProductionCompanies[i].LogoPath, "w185")
	}

	for i := range ent.Networks {
		ent.Networks[i].LogoPath = ImageURL(ent.Networks[i].LogoPath, "w185")
	}

	for i := range ent.Seasons {
		ent.Seasons[i].PosterPath = ImageURL(ent.Seasons[i].PosterPath, "w342")
		ent.Seasons[i].AirDate = FixDate(ent.Seasons[i].AirDate)
	}

	for i := range ent.CreatedBy {
		ent.CreatedBy[i].ProfilePath = ImageURL(ent.CreatedBy[i].ProfilePath, "h632")
	}

	if len(ent.Genres) == 0 && len(ent.GenresIds) > 0 {
		for _, gi := range ent.GenresIds {
			for _, g := range genres {
				if g.ID == gi {
					ent.Genres = append(ent.Genres, g)
					break
				}
			}
		}
	}
	if len(ent.Genres) > 0 && len(ent.GenresIds) == 0 {
		for _, g := range ent.Genres {
			ent.GenresIds = append(ent.GenresIds, g.ID)
		}
	}
}
