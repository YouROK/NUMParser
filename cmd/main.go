package main

import (
	"NUMParser/config"
	"NUMParser/db"
	"NUMParser/db/models"
	"NUMParser/movies/tmdb"
	"NUMParser/parser"
	"NUMParser/releases"
	"NUMParser/web"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type args struct {
	Port  string `arg:"-p" help:"web server port, default 38888"`
	Proxy string `arg:"--proxy" help:"proxy for rutor, http://user:password@ip:port"`
}

var params args

func main() {
	arg.MustParse(&params)

	if params.Port == "" {
		params.Port = "38888"
	}

	if params.Proxy != "" {
		_, err := url.Parse(params.Proxy)
		if err != nil {
			log.Println("Error parse proxy host:", err)
		} else {
			config.ProxyHost = params.Proxy
		}
	}

	db.Init()
	tmdb.Init()

	web.Start(params.Port)

	scanReleases()
	scanMoviesYears()
	web.SetStaticReleases()

	log.Println("Start timer")
	gocron.Every(3).Hours().From(calcTime()).Do(scanReleases)
	gocron.Every(1).Day().At("2:30").Do(scanMoviesYears)
	<-gocron.Start()

}

func scanReleases() {
	getDbInfo()
	rutorParser := parser.NewRutor()
	rutorParser.Parse()

	releases.GetNewMovies()
	releases.GetFourKMovies()
	releases.GetNewTVs()
	releases.GetNewCartoons()
	releases.GetNewCartoonsTV()
	db.SaveAll()
	copy()
}

func scanMoviesYears() {
	releases.GetLegends()
	for y := 1980; y <= time.Now().Year(); y++ {
		releases.GetNewMoviesYear(y)
	}
	db.SaveAll()
	copy()
}

func copy() {
	dir := filepath.Dir(os.Args[0])
	logOut, err := exec.Command("/bin/sh", filepath.Join(dir, "copy.sh")).CombinedOutput()
	if err != nil {
		log.Println("Error copy releases:", err)
	}
	output := string(logOut)
	log.Println(output)
}

func calcTime() *time.Time {
	//2 5 8 11 14 17 20 23
	hour := time.Now().Hour()
	t := time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		0, 0, 0, 0, time.Local)
	if hour < 2 {
		t = t.Add(2 * time.Hour)
	} else if hour < 5 {
		t = t.Add(5 * time.Hour)
	} else if hour < 11 {
		t = t.Add(11 * time.Hour)
	} else if hour < 14 {
		t = t.Add(14 * time.Hour)
	} else if hour < 17 {
		t = t.Add(17 * time.Hour)
	} else if hour < 20 {
		t = t.Add(20 * time.Hour)
	} else if hour < 23 {
		t = t.Add(23 * time.Hour)
	} else if hour >= 23 {
		t = t.Add(26 * time.Hour)
	}
	return &t
}

func getDbInfo() {
	listTorr := db.GetTorrs()
	wIMDB := 0
	nIMDB := 0

	cMovie := 0
	cSeries := 0
	cDocMovie := 0
	cDocSeries := 0
	cCartoonMovie := 0
	cCartoonSeries := 0
	cTVShow := 0
	cAnime := 0

	for _, d := range listTorr {
		if d.IMDBID != "" {
			wIMDB++
		} else {
			nIMDB++
		}
		switch d.Categories {
		case models.CatMovie:
			cMovie++
		case models.CatSeries:
			cSeries++
		case models.CatDocMovie:
			cDocMovie++
		case models.CatDocSeries:
			cDocSeries++
		case models.CatCartoonMovie:
			cCartoonMovie++
		case models.CatCartoonSeries:
			cCartoonSeries++
		case models.CatTVShow:
			cTVShow++
		case models.CatAnime:
			cAnime++
		}
	}

	listTmdb := db.GetTMDBDetails()

	indxs := db.GetIndexes()

	fmt.Println("Movies:", cMovie)
	fmt.Println("Serials:", cSeries)
	fmt.Println("Doc Movies:", cDocMovie)
	fmt.Println("Doc Serials:", cDocSeries)
	fmt.Println("Cartoons:", cCartoonMovie)
	fmt.Println("Cartoons Serial:", cCartoonSeries)
	fmt.Println("TV  Show:", cTVShow)
	fmt.Println("Animes:", cAnime)

	fmt.Println("Torrents with IMDB:", wIMDB)
	fmt.Println("Torrents without IMDB:", nIMDB)
	fmt.Println("Torrents:", len(listTorr))
	fmt.Println("Tmdb details:", len(listTmdb))
	fmt.Println("Indexes:", len(indxs))
}
