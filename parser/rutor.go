package parser

import (
	"NUMParser/db"
	"NUMParser/db/models"
	"NUMParser/tasker"
	"NUMParser/utils"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func loadRutorHost() []string {
	dir := filepath.Dir(os.Args[0])
	name := filepath.Join(dir, "rutor_host.txt")
	buf, err := os.ReadFile(name)
	if err == nil {
		list := strings.Split(string(buf), "\n")
		var ret []string
		for _, l := range list {
			if strings.HasPrefix(l, "http") {
				ret = append(ret, strings.TrimSpace(l))
			}
		}
		return ret
	}
	return nil
}

var (
	hostsPos = 0
	mhost    sync.Mutex
)

type RutorParser struct {
	mu      sync.Mutex
	tasker  *tasker.Tasker
	isParse bool
}

func NewRutor() *RutorParser {
	rt := new(RutorParser)
	rt.tasker = tasker.New(10)
	return rt
}

type parseLink struct {
	Host string
	Link string
	Cat  string
}

func getHost() string {
	mhost.Lock()
	defer mhost.Unlock()
	hosts := loadRutorHost()
	if len(hosts) == 0 {
		return "http://rutor.info"
	}
	host := hosts[hostsPos]
	hostsPos++
	if hostsPos >= len(hosts) {
		hostsPos = 0
	}
	return host
}

func (self *RutorParser) Parse() {
	self.mu.Lock()
	if self.isParse {
		self.mu.Unlock()
		return
	}
	self.isParse = true
	defer func() { self.isParse = false }()
	self.mu.Unlock()
	pages := self.readCategories()
	fullScan := len(db.GetTorrs()) < 10000

	for cat, pgs := range pages {
		for i := 0; i < pgs; i++ {
			page := strconv.Itoa(i)
			host := getHost()
			link := host + "/browse/" + page + "/" + cat + "/0/0"
			pl := parseLink{
				Host: host,
				Link: link,
				Cat:  cat,
			}
			if fullScan {
				self.tasker.Add(func() {
					list := self.parsePage(pl)
					for _, d := range list {
						db.AddTorr(d)
					}
				})
			} else {
				list := self.parsePage(pl)
				dbtorrs := db.GetTorrs()
				finds := 0
				for _, d := range list {
					for _, dbtorr := range dbtorrs {
						if d.Hash == dbtorr.Hash {
							finds++
						}
					}
					db.AddTorr(d)
				}
				if finds == len(list) {
					break
				}
			}
		}
	}

	if fullScan {
		self.tasker.Run()
		self.tasker.Wait()
		db.SaveTorrs()
	} else {
		// Если сканирование не дошло до конца смотрим в конце торренты
		for cat, pgs := range pages {
			for i := pgs; i > 0; i++ {
				page := strconv.Itoa(i)
				host := getHost()
				link := host + "/browse/" + page + "/" + cat + "/0/0"
				pl := parseLink{
					Host: host,
					Link: link,
					Cat:  cat,
				}
				list := self.parsePage(pl)
				dbtorrs := db.GetTorrs()
				finds := 0
				for _, d := range list {
					for _, dbtorr := range dbtorrs {
						if d.Hash == dbtorr.Hash {
							finds++
						}
					}
					db.AddTorr(d)
				}
				if finds == len(list) {
					break
				}
			}
		}
		db.SaveTorrs()
	}
}

// читаем категории и заносим в таски что парсить
func (self *RutorParser) readCategories() map[string]int {
	// 1  - Зарубежные фильмы          	| Фильмы
	// 5  - Наши фильмы                	| Фильмы
	// 4  - Зарубежные сериалы         	| Сериалы
	// 16 - Наши сериалы               	| Сериалы
	// 12 - Научно-популярные фильмы   	| Док. сериалы, Док. фильмы
	// 6  - Телевизор                  	| ТВ Шоу
	// 7  - Мультипликация             	| Мультфильмы, Мультсериалы
	// 10 - Аниме                      	| Аниме
	// 17 - Иностранные релизы         	| UA озвучка
	// 13 - Спорт и Здоровье 			| ТВ Шоу
	// 15 - Юмор						| ТВ Шоу

	log.Println("Read Rutor categories")

	var categories = []string{"1", "5", "4", "16", "12", "6", "7", "10", "17", "13", "15"}
	var pages = map[string]int{}
	var mm sync.Mutex

	utils.PFor(categories, func(i int, cat string) {
		link := getHost() + "/browse/0/" + cat + "/0/0"
		body, err := get(link)
		if err == nil {
			re, err := regexp.Compile("<a href=\"/browse/([0-9]+)/[0-9]+/[0-9]+/[0-9]+\"><b>[0-9]+&nbsp;-&nbsp;[0-9]+</b></a></p>")
			if err != nil {
				log.Fatalf("Error compile regex %v", err)
			}
			matches := re.FindStringSubmatch(body)
			if len(matches) > 1 {
				pgs, err := strconv.Atoi(matches[1])
				if err == nil {
					mm.Lock()
					pages[cat] = pgs
					log.Println("Category readed", link, pgs)
					mm.Unlock()
				}
			}
		}
	})

	return pages
}

func (self *RutorParser) parsePage(pl parseLink) []*models.TorrentDetails {
	// Парсим страницу с торрентами
	body, err := get(pl.Link)
	if err != nil {
		log.Println("Error get page:", err, pl.Link)
		return nil
	}
	if !strings.Contains(body, "<title>rutor.info") {
		log.Println("Not rutor page:", pl.Link)
		return nil
	}
	log.Println("Readed:", pl.Link)
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(body))
	if err != nil {
		log.Println("Error parse page:", err, pl.Link)
		return nil
	}

	var list []*models.TorrentDetails

	doc.Find("div#index").Find("tr").Each(func(_ int, selection *goquery.Selection) {
		if selection.HasClass("backgr") {
			return
		}
		selTd := selection.Find("td")

		itm := new(models.TorrentDetails)
		itm.CreateDate = self.parseDate(node2Text(selTd.Get(0)))
		itm.Title = node2Text(selTd.Get(1))
		self.parseTitle(itm, pl.Cat)
		itm.Magnet = selTd.Get(1).FirstChild.NextSibling.Attr[0].Val
		hash := getHash(itm.Magnet)
		if hash == "" {
			return
		}
		itm.Hash = hash
		linkParam := selTd.Get(1).LastChild.Attr[0].Val
		itm.Link = linkParam
		if len(selTd.Nodes) == 4 {
			itm.Size = node2Text(selTd.Get(2))
			peers := node2Text(selTd.Get(3))
			prarr := strings.Split(peers, "  ")
			if len(prarr) > 1 {
				itm.Seed, _ = strconv.Atoi(prarr[0])
				itm.Peer, _ = strconv.Atoi(prarr[1])
			}
		} else if len(selTd.Nodes) == 5 {
			itm.Size = node2Text(selTd.Get(3))
			peers := node2Text(selTd.Get(4))
			prarr := strings.Split(peers, "  ")
			if len(prarr) > 1 {
				itm.Seed, _ = strconv.Atoi(prarr[0])
				itm.Peer, _ = strconv.Atoi(prarr[1])
			}
		}
		itm.Tracker = "Rutor"
		list = append(list, itm)
	})
	return list
}

func (self *RutorParser) parseTitle(td *models.TorrentDetails, cat string) {
	td.Title = strings.ReplaceAll(td.Title, "&amp;", "&")
	re, err := regexp.Compile("(.+)\\((.+)\\)(.+)")
	if err != nil {
		log.Fatalf("Error parse torrent name:", err)
	}
	matches := re.FindStringSubmatch(td.Title)
	re, err = regexp.Compile("\\[.*?\\]")
	if len(matches) > 2 {
		yrs := strings.TrimSpace(matches[2])
		if strings.Contains(yrs, "-") {
			arr := strings.Split(yrs, "-")
			if len(arr) > 1 {
				yrs = strings.TrimSpace(arr[0])
			}
		}

		yr, _ := strconv.Atoi(yrs)
		td.Year = yr
		arr := strings.Split(matches[1], "/")
		if len(arr) == 1 {
			td.Name = arr[0]
		} else if len(arr) > 1 {
			td.Name = arr[0]
			td.Names = arr[1:]
		}
		if err == nil {
			if strings.Contains(td.Title, "З/Л/О 94") {
				td.Name = "З/Л/О 94"
				td.Names = []string{"V/H/S/94"}
			} else {
				td.Name = strings.TrimSpace(re.ReplaceAllString(td.Name, ""))
				for i := range td.Names {
					td.Names[i] = strings.TrimSpace(re.ReplaceAllString(td.Names[i], ""))
				}
			}
		}
		if len(matches) > 3 {
			vq := ParseVQuality(strings.TrimSpace(matches[3]))
			aq := ParseAQuality(strings.TrimSpace(matches[3]))
			td.VideoQuality = vq
			td.AudioQuality = aq
		}
	}

	title := td.Title
	if len(matches) > 0 {
		title = matches[1]
	}

	switch {
	case cat == "1", cat == "5", cat == "17":
		td.Categories = models.CatMovie
	case cat == "4", cat == "16":
		td.Categories = models.CatSeries
	case cat == "12":
		if re.MatchString(title) {
			td.Categories = models.CatDocSeries
		} else {
			td.Categories = models.CatDocMovie
		}
	case cat == "6", cat == "13", cat == "15":
		td.Categories = models.CatTVShow
	case cat == "7":
		if re.MatchString(title) {
			td.Categories = models.CatCartoonSeries
		} else {
			td.Categories = models.CatCartoonMovie
		}
	case cat == "10":
		td.Categories = models.CatAnime
	}
}

func (self *RutorParser) parseDate(date string) time.Time {
	var rutorMonth = map[string]int{
		"Янв": 1, "Фев": 2, "Мар": 3,
		"Апр": 4, "Май": 5, "Июн": 6,
		"Июл": 7, "Авг": 8, "Сен": 9,
		"Окт": 10, "Ноя": 11, "Дек": 12,
	}

	darr := strings.Split(date, " ")
	if len(darr) != 3 {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.Now().Location())
	}

	day, _ := strconv.Atoi(darr[0])
	month, _ := rutorMonth[darr[1]]
	year, _ := strconv.Atoi("20" + darr[2])

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location())
}

func GetBodyLink(torr *models.TorrentDetails) string {
	body, err := get(getHost() + torr.Link)
	if err != nil {
		log.Println("Error get page:", err, torr.Link)
		return ""
	}
	return body
}

func ParseVQuality(params string) int {
	info := clear(strings.ToLower(params))
	info = strings.ReplaceAll(info, "вdrip", "bdrip")

	// check uhd bdremux 2160
	if strings.Contains(info, "2160") && (strings.Contains(info, "bdremux") || strings.Contains(info, "bluray")) {
		if strings.Contains(info, "dolby vision") {
			return models.Q_UHD_BDREMUX_DV
		} else if strings.Contains(info, "hdr") {
			return models.Q_UHD_BDREMUX_HDR
		} else {
			return models.Q_UHD_BDREMUX_SDR
		}
	}
	// check bdrip hevc 2160
	if strings.Contains(info, "2160") && strings.Contains(info, "bdrip") {
		if strings.Contains(info, "dolby vision") {
			return models.Q_BDRIP_DV_2160
		} else if strings.Contains(info, "hdr") {
			return models.Q_BDRIP_HDR_2160
		} else {
			return models.Q_BDRIP_SDR_2160
		}
	}
	// check webdl 2160
	if strings.Contains(info, "2160") && (strings.Contains(info, "webdl") || (strings.Contains(info, "webrip"))) {
		if strings.Contains(info, "dolby vision") {
			return models.Q_WEBDL_DV_2160
		} else if strings.Contains(info, "hdr") {
			return models.Q_WEBDL_HDR_2160
		} else {
			return models.Q_WEBDL_SDR_2160
		}
	}
	// check bdremux 1080
	if strings.Contains(info, "1080") && (strings.Contains(info, "remux") || strings.Contains(info, "bluray")) {
		return models.Q_BDREMUX_1080
	}
	// check bdrip hevc 1080
	if strings.Contains(info, "1080") && strings.Contains(info, "bdrip") && strings.Contains(info, "hevc") {
		return models.Q_BDRIP_HEVC_1080
	}
	// check bdrip 1080
	if strings.Contains(info, "1080") && strings.Contains(info, "bdrip") {
		return models.Q_BDRIP_1080
	}
	// check webdl 1080
	if strings.Contains(info, "1080") && (strings.Contains(info, "webdl") || strings.Contains(info, "webrip") || strings.Contains(info, "hdrip") || strings.Contains(info, "hybrid")) {
		return models.Q_WEBDL_1080
	}
	// check bdrip hevc 720
	if strings.Contains(info, "720") && strings.Contains(info, "bdrip") && strings.Contains(info, "hevc") {
		return models.Q_BDRIP_HEVC_720
	}
	// check bdrip 720
	if strings.Contains(info, "720") && strings.Contains(info, "bdrip") {
		return models.Q_BDRIP_720
	}
	// check webdl 720
	if strings.Contains(info, "720") && (strings.Contains(info, "webdl") || strings.Contains(info, "webrip") || strings.Contains(info, "dvd") || strings.Contains(info, "hdrip")) {
		return models.Q_WEBDL_720
	}
	return models.Q_LOWER
}

func ParseAQuality(params string) int {
	arr := strings.Split(params, "|")
	qualities := []int{}
	for _, name := range arr {
		name = clear(name)
		for _, qn := range Q_Lic_Names {
			if strings.Contains(name, clear(qn)) {
				qualities = append(qualities, models.Q_LICENSE)
			}
		}
		// Ищем проф студию
		for _, qn := range Q_P_Names {
			// одно слово ищем по словам
			if !strings.Contains(qn, " ") {
				arr := strings.Split(clear(name), " ")
				for _, s := range arr {
					if s == clear(qn) {
						qualities = append(qualities, models.Q_PS)
					}
				}
			} else {
				if strings.Contains(name, clear(qn)) {
					qualities = append(qualities, models.Q_PS)
				}
			}
		}
		// Ищем любительскую студию
		for _, qn := range Q_L_Names {
			if strings.Contains(name, clear(qn)) {
				qualities = append(qualities, models.Q_LS)
			}
		}
		// Ищем в названии упоминания озвучки
		wrd := strings.Split(name, " ")
		for _, w := range wrd {
			w = strings.TrimSpace(w)
			if w == "d" {
				qualities = append(qualities, models.Q_D)
			} else if w == "p" {
				qualities = append(qualities, models.Q_P)
			} else if w == "p2" {
				qualities = append(qualities, models.Q_P2)
			} else if w == "p1" {
				qualities = append(qualities, models.Q_P1)
			} else if w == "l" {
				qualities = append(qualities, models.Q_L)
			} else if w == "l2" {
				qualities = append(qualities, models.Q_L2)
			} else if w == "l1" {
				qualities = append(qualities, models.Q_L1)
			} else if w == "a" {
				qualities = append(qualities, models.Q_A)
			}
		}
	}

	if len(qualities) > 0 {
		max := models.Q_UNKNOWN
		for _, q := range qualities {
			if q > max {
				max = q
			}
		}
		return max
	}

	return models.Q_UNKNOWN
}

func clear(txt string) string {
	ret := ""
	txt = strings.ToLower(txt)
	for _, r := range txt {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'а' && r <= 'я') || r == 'ё' || r == ' ' {
			ret = ret + string(r)
		}
	}
	return ret
}

var Q_Lic_Names = []string{
	"лицензия",
	"itunes",
	"netflix",
}

var Q_P_Names = []string{
	"100ТВ",
	"2х2",
	"Agatha Studdio",
	"AlexFilm",
	"Amedia",
	"NovaFilm",
	"Novamedia",
	"AMS",
	"ARS-studio",
	"Astana TV",
	"AzOnFilm",
	"AXN Sci-Fi",
	"CDV",
	"CGInfo",
	"CP Digital",
	"Disney",
	"DniproFilm",
	"DVDXpert",
	"Elrom",
	"Filiza Studio",
	"Flarrow Films",
	"FocusX",
	"FocusStudio",
	"FOXCrime",
	"FoxLife",
	"Gears Media",
	"Good People",
	"HDrezka Studio",
	"IdeaFilm",
	"IVI",
	"Jaskier",
	"Kansai Studio",
	"LostFilm",
	"MC Entertaiment",
	"Mega-Anime",
	"MTV",
	"Neoclassica",
	"NewComers",
	"NewStudio",
	"Nickelodeon",
	"NovaFilm",
	"NovaMedia",
	"Ozz",
	"Paramount",
	"Profix Media",
	"Rattlebox",
	"SDI Media",
	"Sony Sci-Fi",
	"Superbit",
	"TUMBLER Studio",
	"TVShows",
	"FilmsClub",
	"Tycoon",
	"Universal",
	"ViruseProject",
	"WestVideo",
	"Арена",
	"Арк-ТВ",
	"Воротилин",
	"Домашний",
	"ДТВ",
	"ДубльPR studio",
	"Екатеринбург-Арт",
	"Инис",
	"Лексикон",
	"Киномания",
	"Кипарис",
	"Кириллица",
	"Кравец",
	"Кубик в Кубе",
	"Кураж Бомбей",
	"Невафильм",
	"Новый канал",
	"НТВ",
	"НТВ+",
	"Омикрон",
	"ОРТ",
	"Парадиз-ВС",
	"Первый канал",
	"Петербург 5 канал",
	"Пифагор",
	"Позитив-Мультимедиа",
	"Премьер Видео Фильм",
	"РЕН-ТВ",
	"С.Р.И.",
	"Специальное Российское Издание",
	"СВ-Студия",
	"СоюзВидео",
	"студия «Велес»",
	"студия «Нота»",
	"студия «СВ Дубль»",
	"ТВ3",
	"ТВЦ",
	"ТНТ",
	"ТОО Прим",
	"Хабар",
	"Pazl Voice",
}

var Q_L_Names = []string{
	"Albion Studio",
	"Alternative Production",
	"AniDub",
	"AniFilm",
	"AniLibria",
	"Anilife Project",
	"AniMedia",
	"AnimeReactor",
	"AnimeVost",
	"AniPlay",
	"AniStar",
	"ApofysTeam",
	"Baibako",
	"BraveSound",
	"CACTUS TEAM",
	"СoldFilm",
	"DexterTV",
	"DreamRecords",
	"Eleonor Film",
	"E-Production",
	"Etvox Film",
	"Filiza Studio",
	"Flux-Team",
	"F-TRAIN",
	"GladiolusTV",
	"GostFilm",
	"Gramalant",
	"GREEN TEA",
	"GSGroup",
	"HamsterStudio",
	"ICG",
	"Jetvis Studio",
	"Jimmy J",
	"LevshaFilm",
	"LugaDUB",
	"LE-Production",
	"Mallorn Studio",
	"MYDIMKA",
	"Naruto-Hokage",
	"NetLab Anima Group",
	"NewStation",
	"NikiStudio Records",
	"OmskBird",
	"OneFilm",
	"OpenDub",
	"Padabajour",
	"ParadoX",
	"RG.Paravozik",
	"Shiza Project",
	"SkyeFilmTV",
	"STEPonee",
	"StopGame",
	"Sunny-Films",
	"To4kaTV",
	"VictoryFilms",
	"Web_Money",
	"ZM-SHOW",
	"Несмертельное оружие",
	"Причудики",
	"Райдо",
	"Синема-УС",
	"Студия Пиратского Дубляжа",
	"Сладкая парочка",
	"Частная Студия",
}
