package api

import (
	"sync"
	"time"
)

const (
	CategoryGeneral = "general"
	CategoryAnime   = "anime"
	CategoryPeople  = "people"
	PuritySFW       = "sfw"
	PuritySketchy   = "sketchy"
	PurityNSFW      = "nsfw"
	SortDate        = "date_added"
	SortRelevance   = "relevance"
	SortRandom      = "random"
	SortViews       = "views"
	SortFavorites   = "favorites"
	SortTopList     = "toplist"
	Range1day       = 24 * time.Hour
	Range3day       = 3 * Range1day
	Range1week      = 7 * Range1day
	Range1month     = 30 * Range1day
	Range3months    = 3 * Range1month
	Range6months    = 6 * Range1month
	Range1year      = 12 * Range1month
	RatioLandscape  = "landscape"
	RatioPortrait   = "portrait"
	ThumbSmall      = "small"
	ThumbLarge      = "large"
	ThumbOriginal   = "original"
)

type Uploader struct {
	Username string `json:"username"`
	Group    string `json:"group"`
}

type Thumbs struct {
	Small    string `json:"small"`
	Large    string `json:"large"`
	Original string `json:"original"`
}

type Tag struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Alias     string `json:"alias"`
	Category  string `json:"category"`
	Purity    string `json:"purity"`
	CreatedAt string `json:"created_at"`
}

type Wallpaper struct {
	Id         string   `json:"id"`
	Url        string   `json:"url"`
	ShortUrl   string   `json:"short_url"`
	Source     string   `json:"source"`
	Uploader   Uploader `json:"uploader"`
	Views      int      `json:"views"`
	Favorites  int      `json:"favorites"`
	Category   string   `json:"category"`
	Purity     string   `json:"purity"`
	DimensionX int      `json:"dimension_x"`
	DimensionY int      `json:"dimension_y"`
	Resolution string   `json:"resolution"`
	Ratio      string   `json:"ratio"`
	FileSize   int      `json:"file_size"`
	FileType   string   `json:"file_type"`
	Path       string   `json:"path"`
	Colors     []string `json:"colors"`
	Thumbs     Thumbs   `json:"thumbs"`
	Tags       []Tag    `json:"tags"`
	CreatedAt  string   `json:"created_at"`
}

type QueryCond struct {
	Query       string
	Categories  string `validate:"oneof=000 001 010 011 100 101 110 111"`
	Purity      string `validate:"oneof=000 001 010 011 100 101 110 111"`
	Sorting     string `validate:"oneof=date_added relevance random views favorites toplist hot"`
	Order       string `validate:"oneof=desc asc"`
	TopRange    string `validate:"oneof=1d 3d 1w 1M 3M 6M 1y"`
	AtLeast     string
	Resolutions string
	Ratios      string
	Colors      string
	Seed        string
	Page        int
}

type QueryResult struct {
	api   *API
	cond  *QueryCond
	pIdx  int
	pLock sync.Mutex

	pickLoop bool
	pickRand sync.Once

	Data []*Wallpaper `json:"data"`
	Meta struct {
		CurrentPage int `json:"current_page"`
		LastPage    int `json:"last_page"`
		Total       int `json:"total"`
	} `json:"meta"`
}
