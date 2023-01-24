package models

import "time"

type KPDetail struct {
	KinopoiskID                int         `json:"kinopoiskId,omitempty"`
	ImdbID                     string      `json:"imdbId,omitempty"`
	NameRu                     string      `json:"nameRu,omitempty"`
	NameEn                     string      `json:"nameEn,omitempty"`
	NameOriginal               string      `json:"nameOriginal,omitempty"`
	PosterURL                  string      `json:"posterUrl,omitempty"`
	PosterURLPreview           string      `json:"posterUrlPreview,omitempty"`
	CoverURL                   string      `json:"coverUrl,omitempty"`
	LogoURL                    string      `json:"logoUrl,omitempty"`
	ReviewsCount               int         `json:"reviewsCount,omitempty"`
	RatingGoodReview           float64     `json:"ratingGoodReview,omitempty"`
	RatingGoodReviewVoteCount  int         `json:"ratingGoodReviewVoteCount,omitempty"`
	RatingKinopoisk            float64     `json:"ratingKinopoisk,omitempty"`
	RatingKinopoiskVoteCount   int         `json:"ratingKinopoiskVoteCount,omitempty"`
	RatingImdb                 float64     `json:"ratingImdb,omitempty"`
	RatingImdbVoteCount        int         `json:"ratingImdbVoteCount,omitempty"`
	RatingFilmCritics          float64     `json:"ratingFilmCritics,omitempty"`
	RatingFilmCriticsVoteCount int         `json:"ratingFilmCriticsVoteCount,omitempty"`
	RatingAwait                float64     `json:"ratingAwait,omitempty"`
	RatingAwaitCount           int         `json:"ratingAwaitCount,omitempty"`
	RatingRfCritics            float64     `json:"ratingRfCritics,omitempty"`
	RatingRfCriticsVoteCount   int         `json:"ratingRfCriticsVoteCount,omitempty"`
	WebURL                     string      `json:"webUrl,omitempty"`
	Year                       int         `json:"year,omitempty"`
	FilmLength                 int         `json:"filmLength,omitempty"`
	Slogan                     string      `json:"slogan,omitempty"`
	Description                string      `json:"description,omitempty"`
	ShortDescription           string      `json:"shortDescription,omitempty"`
	EditorAnnotation           string      `json:"editorAnnotation,omitempty"`
	IsTicketsAvailable         bool        `json:"isTicketsAvailable,omitempty"`
	ProductionStatus           string      `json:"productionStatus,omitempty"`
	Type                       string      `json:"type,omitempty"`
	RatingMpaa                 string      `json:"ratingMpaa,omitempty"`
	RatingAgeLimits            string      `json:"ratingAgeLimits,omitempty"`
	HasImax                    bool        `json:"hasImax,omitempty"`
	Has3D                      bool        `json:"has3D,omitempty"`
	Countries                  []Countries `json:"countries,omitempty"`
	Genres                     []Genres    `json:"genres,omitempty"`
	StartYear                  int         `json:"startYear,omitempty"`
	EndYear                    int         `json:"endYear,omitempty"`
	Serial                     bool        `json:"serial,omitempty"`
	ShortFilm                  bool        `json:"shortFilm,omitempty"`
	Completed                  bool        `json:"completed,omitempty"`
	UpdateDate                 time.Time   `json:"update_date,omitempty"`
}

type KPSearch struct {
	Keyword                string  `json:"keyword,omitempty"`
	PagesCount             int     `json:"pagesCount,omitempty"`
	Films                  []Films `json:"films,omitempty"`
	SearchFilmsCountResult int     `json:"searchFilmsCountResult,omitempty"`
}

type Films struct {
	FilmID           int         `json:"filmId,omitempty"`
	NameRu           string      `json:"nameRu,omitempty"`
	NameEn           string      `json:"nameEn,omitempty"`
	Type             string      `json:"type,omitempty"`
	Year             string      `json:"year,omitempty"`
	Description      string      `json:"description,omitempty"`
	FilmLength       string      `json:"filmLength,omitempty"`
	Countries        []Countries `json:"countries,omitempty"`
	Genres           []Genres    `json:"genres,omitempty"`
	Rating           string      `json:"rating,omitempty"`
	RatingVoteCount  int         `json:"ratingVoteCount,omitempty"`
	PosterURL        string      `json:"posterUrl,omitempty"`
	PosterURLPreview string      `json:"posterUrlPreview,omitempty"`
}

type Countries struct {
	Country string `json:"country,omitempty"`
}

type Genres struct {
	Genre string `json:"genre,omitempty"`
}
