package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GoGoAnimeApi struct {
	BaseURL   string
	SearchURL string
}

func NewGoGoAnimeApi() GoGoAnimeApi {
	return GoGoAnimeApi{
		SearchURL: "https://ajax.gogo-load.com/site/loadAjaxSearch",
		BaseURL:   "https://www1.gogoanime.cm",
	}
}

func (gogoApi GoGoAnimeApi) GetEpisode(result SearchResult, number uint64) (Episode, error) {
	detailPageURLSplit := strings.Split(result.DetailPageURL, "/")

	animeTagName := detailPageURLSplit[len(detailPageURLSplit)-1]
	res, err := http.Get(fmt.Sprintf("%s/%s-episode-%d", gogoApi.BaseURL, animeTagName, number))
	if err != nil {
		return Episode{}, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Episode{}, err

	}

	streamListURL, _ := doc.Find(".dowloads").First().Find("a").Attr("href")
	streamListHTML, err := http.Get(streamListURL)
	if err != nil {
		return Episode{}, err
	}

	streamListDoc, err := goquery.NewDocumentFromReader(streamListHTML.Body)
	if err != nil {
		return Episode{}, err
	}

	var streams []StreamSource
	streamListDoc.Find(".mirror_link").First().Find(".dowload").Each(func(i int, s *goquery.Selection) {
		anchor := s.Find("a")
		streamType := strings.TrimSpace(strings.Split(anchor.Text(), "\n")[1])
		href, _ := anchor.Attr("href")
		streams = append(streams, StreamSource{URL: href, Type: streamType, Origin: gogoApi.Tag()})
	})

	return Episode{
		Number:        number,
		StreamSources: streams,
	}, nil
}

func (GoGoAnimeApi) Tag() AnimeApiTag {
	return GoGoAnime
}

func (gogoApi GoGoAnimeApi) GetDetail(result SearchResult) (Detail, error) {
	pageHTML, err := http.Get(result.DetailPageURL)
	if err != nil {
		return Detail{}, err
	}

	pageDoc, err := goquery.NewDocumentFromReader(pageHTML.Body)
	if err != nil {
		return Detail{}, err
	}

	Type := pageDoc.Clone().Find(".type").First().Find("a").Text()
	synopsisDoc := pageDoc.Clone().Find(".type").Eq(1)
	synopsisDoc.Find("span").Remove()
	synopsis := synopsisDoc.Text()

	episodesString, _ := pageDoc.Clone().Find("#episode_page").Find("a").Attr("ep_end")
	episodes, err := strconv.ParseUint(episodesString, 10, 64)
	if err != nil {
		return Detail{}, err
	}

	return Detail{
		Type:     Type,
		Synopsis: synopsis,
		Episodes: episodes,
	}, nil
}

func (gogoApi GoGoAnimeApi) Search(name string) ([]SearchResult, error) {
	resp, err := http.Get(gogoApi.SearchURL + "?keyword=" + name)
	if err != nil {
		return []SearchResult{}, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return []SearchResult{}, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return []SearchResult{}, err
	}

	htmlString, ok := data["content"].(string)
	if !ok {
		return []SearchResult{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf("<html>%s</html>", htmlString)))
	if err != nil {
		return []SearchResult{}, err
	}

	var results []SearchResult
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		name := s.Text()
		imgStyle, _ := s.Find("div").Attr("style")
		r := regexp.MustCompile(`url\("(.*)"\)`)
		imgSrc := r.FindStringSubmatch(imgStyle)[1]
		pageURL, _ := s.Attr("href")

		result := SearchResult{
			Title:         name,
			ImageSrc:      imgSrc,
			DetailPageURL: gogoApi.BaseURL + "/" + pageURL,
		}

		results = append(results, result)
	})

	return results, nil
}
