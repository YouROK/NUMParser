package releases

type TmdbId struct {
	Id          int64      `json:"id,omitempty"`
	MediaType   string     `json:"media_type,omitempty"`
	Torrent     []*Torrent `json:"torrent,omitempty"`
	GenreIds    []int      `json:"genre_ids,omitempty"`
	VoteAverage float64    `json:"vote_average,omitempty"`
	VoteCount   int        `json:"vote_count,omitempty"`
	Countries   []string   `json:"countries,omitempty"`
	ReleaseDate string     `json:"release_date,omitempty"`
}

type ReleasesID struct {
	Date  string    `json:"date,omitempty"`
	Time  string    `json:"time,omitempty"`
	Items []*TmdbId `json:"items,omitempty"`
}

type Torrent struct {
	Name     string `json:"name,omitempty"`
	Date     string `json:"date,omitempty"`
	Magnet   string `json:"magnet,omitempty"`
	Size     string `json:"size,omitempty"`
	Upload   string `json:"upload,omitempty"`
	Download string `json:"download,omitempty"`
	Source   string `json:"source,omitempty"`
	Link     string `json:"link,omitempty"`
	Quality  int    `json:"quality,omitempty"`
	Voice    int    `json:"voice,omitempty"`
}

type CollectionId struct {
	Name         string    `json:"name"`
	Overview     string    `json:"overview"`
	Parts        []*TmdbId `json:"parts"`
	PosterPath   string    `json:"poster_path"`
	BackdropPath string    `json:"backdrop_path"`
}
