package models

import "time"

type Entity struct {
	Adult               bool                 `json:"adult"`
	BackdropPath        string               `json:"backdrop_path"`
	BelongsToCollection *BelongsToCollection `json:"belongs_to_collection"`
	Budget              int64                `json:"budget"`
	GenresIds           []int                `json:"genre_ids"`
	Genres              []*Genre             `json:"genres"`
	Homepage            string               `json:"homepage"`
	ID                  int64                `json:"id"`
	ImdbID              string               `json:"imdb_id"`
	OriginalLanguage    string               `json:"original_language"`
	OriginalTitle       string               `json:"original_title"`
	Overview            string               `json:"overview"`
	Popularity          float64              `json:"popularity"`
	PosterPath          string               `json:"poster_path"`
	ProductionCompanies []*ProductionCompany `json:"production_companies"`
	ProductionCountries []*ProductionCountry `json:"production_countries"`
	ReleaseDate         string               `json:"release_date"`
	Revenue             int64                `json:"revenue"`
	Runtime             int                  `json:"runtime"`
	SpokenLanguages     []*SpokenLanguage    `json:"spoken_languages"`
	Status              string               `json:"status"`
	Tagline             string               `json:"tagline"`
	Title               string               `json:"title"`
	Video               bool                 `json:"video"`
	VoteAverage         float64              `json:"vote_average"`
	VoteCount           int                  `json:"vote_count"`
	Titles              []string             `json:"titles"`

	// tv
	CreatedBy        []*CreatedBy `json:"created_by"`
	EpisodeRunTime   []int        `json:"episode_run_time"`
	FirstAirDate     string       `json:"first_air_date"`
	InProduction     bool         `json:"in_production"`
	Languages        []string     `json:"languages"`
	LastAirDate      string       `json:"last_air_date"`
	Name             string       `json:"name"`
	Networks         []*Network   `json:"networks"`
	NumberOfEpisodes int          `json:"number_of_episodes"`
	NumberOfSeasons  int          `json:"number_of_seasons"`
	OriginCountry    []string     `json:"origin_country"`
	OriginalName     string       `json:"original_name"`
	Seasons          []*Season    `json:"seasons"`
	Type             string       `json:"type"`
	Images           *Images      `json:"images,omitempty"`

	//multi
	Year         string `json:"year"`
	Character    string `json:"character"`
	MediaType    string `json:"media_type"`
	CreditID     string `json:"credit_id"`
	EpisodeCount int    `json:"episode_count,omitempty"`

	UpdateDate time.Time `json:"update_date,omitempty"`
	torrent    *TorrentDetails
}

func (e *Entity) GetTorrent() *TorrentDetails {
	return e.torrent
}

func (e *Entity) SetTorrent(t *TorrentDetails) {
	e.torrent = t
}

type AlternativeTitles struct {
	ID     int       `json:"id"`
	Titles []*Titles `json:"titles"`
}
type Titles struct {
	Iso31661 string `json:"iso_3166_1"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type BelongsToCollection struct {
	Id       int    `json:"id"`
	Name     string `json:"name,omitempty"`
	Poster   string `json:"poster_path,omitempty"`
	Backdrop string `json:"backdrop_path,omitempty"`
}

type ProductionCompany struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCountry struct {
	Iso31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

type SpokenLanguage struct {
	Iso6391 string `json:"iso_639_1"`
	Name    string `json:"name"`
}

type CreatedBy struct {
	ID          int    `json:"id"`
	CreditID    string `json:"credit_id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	ProfilePath string `json:"profile_path"`
}

type Network struct {
	Name          string `json:"name"`
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

type Season struct {
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

type EntityRequest struct {
	Page         int       `json:"page"`
	Results      []*Entity `json:"results"`
	TotalPages   int       `json:"total_pages"`
	TotalResults int       `json:"total_results"`
}

type CollectionRequest struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Overview     string    `json:"overview"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath string    `json:"backdrop_path"`
	Parts        []*Entity `json:"parts,omitempty"`
}

type Images struct {
	Backdrops []struct {
		AspectRatio float64 `json:"aspect_ratio"`
		FilePath    string  `json:"file_path"`
		Height      int     `json:"height"`
		Iso6391     string  `json:"iso_639_1"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		Width       int     `json:"width"`
	} `json:"backdrops"`
	Posters []struct {
		AspectRatio float64 `json:"aspect_ratio"`
		FilePath    string  `json:"file_path"`
		Height      int     `json:"height"`
		Iso6391     string  `json:"iso_639_1"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		Width       int     `json:"width"`
	} `json:"posters"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GenresRequest struct {
	Genres []*Genre `json:"genres"`
}

type Configuration struct {
	Images struct {
		BaseURL       string   `json:"base_url"`
		SecureBaseURL string   `json:"secure_base_url"`
		BackdropSizes []string `json:"backdrop_sizes"`
		LogoSizes     []string `json:"logo_sizes"`
		PosterSizes   []string `json:"poster_sizes"`
		ProfileSizes  []string `json:"profile_sizes"`
		StillSizes    []string `json:"still_sizes"`
	}
	ChangeKeys []string `json:"change_keys,omitempty"`
}

type TrailersRequest struct {
	ID      int        `json:"id"`
	Results []*Trailer `json:"results"`
}

type Trailers []*Trailer

type Trailer struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	Site string `json:"site"`
	Key  string `json:"key"`
	Type string `json:"type"`

	Link   string `json:"link"`
	Poster string `json:"poster"`
}

type FindResult struct {
	MovieResults []*Entity `json:"movie_results,omitempty"`
	TVResults    []*Entity `json:"tv_results,omitempty"`
}
