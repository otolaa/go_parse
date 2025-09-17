package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ItemParse struct {
	Title   string  `json:"title"`
	Href    string  `json:"href"`
	Img     string  `json:"img"`
	Episode Episode `json:"episode"`
	Rating  Rating  `json:"rating"`
}

type Episode struct {
	Season int `json:"season"`
	Series int `json:"series"`
}

type Rating struct {
	Kp   float64 `json:"kp"`
	Imdb float64 `json:"imdb"`
}

var Items []ItemParse

// 1 сезон 10 серия
func GetEpisode(a *goquery.Selection) (Episode, error) {
	episode := a.Find("div.th-series").Text()
	epi := strings.TrimSpace(episode)
	if epi == "" {
		return Episode{}, errors.New("not episode")
	}

	episodeArr := strings.Split(epi, " ")
	if len(episodeArr[0]) == 0 || len(episodeArr[2]) == 0 {
		return Episode{}, errors.New("not season, series")
	}

	season, err := strconv.Atoi(episodeArr[0])
	if err != nil {
		return Episode{}, errors.New("not season")
	}

	series, err := strconv.Atoi(episodeArr[2])
	if err != nil {
		return Episode{}, errors.New("not series")
	}

	return Episode{Season: season, Series: series}, nil
}

func GetRating(a *goquery.Selection) Rating {
	rates := a.Find("div.th-desc > div.th-rates")
	kp := rates.Find("div.th-rate.th-rate-kp > span").Text()
	imbd := rates.Find("div.th-rate.th-rate-imdb > span").Text()

	kpNum, err := strconv.ParseFloat(kp, 64)
	if err != nil {
		kpNum = 0
	}

	imbdNum, err := strconv.ParseFloat(imbd, 64)
	if err != nil {
		imbdNum = 0
	}

	return Rating{Kp: kpNum, Imdb: imbdNum}
}

// For each item found, get the title
func GetItem(i int, s *goquery.Selection) {
	a := s.Find("a.th-in.with-mask")
	href, _ := a.Attr("href")
	title := a.Find("div.th-desc > div.th-title").Text()
	img, _ := a.Find("div.th-img > img").Attr("src")

	item := ItemParse{
		Title:  strings.TrimSpace(title),
		Href:   href,
		Img:    img,
		Rating: GetRating(a),
	}

	episode, err := GetEpisode(a)
	if err == nil {
		item.Episode = episode
	}

	Items = append(Items, item)

	p(2, " → ", "[+]", i, href, title, episode)
}

// Request the HTML page.
func GetHtml(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		e := fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New(e)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func GetScrape(url string) []ItemParse {
	doc, err := GetHtml(url)
	if err != nil {
		log.Fatal(err)
	}

	Items = nil

	// Find the review items
	doc.Find("#dle-content > div.th-item").Each(GetItem)

	return Items
}

func main() {
	Url := "https://lordserial.run/zarubezhnye-serialy/"
	UrlPage := Url
	f := "./json/item.json"
	for v := range 55 {
		if v > 0 {
			UrlPage = fmt.Sprintf("%spage/%d/", Url, v+1)
		}

		FilmItems := GetScrape(UrlPage)

		p(3, " ~ ", "[+]", v+1, UrlPage, len(FilmItems))

		// append json
		appendJson(FilmItems, f)
	}
}
