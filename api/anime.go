package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	Id          string
	Title       string
	Description string
	Type        string
	ImageSrc    string
	PageURL     string
}

type Detail struct {
	Synopsis string
	Type     string
	Episodes uint64
	MalURL   string
}

type StreamSource struct {
	URL    string
	Type   string
	Origin string
}

type Episode struct {
	Number        uint64
	StreamSources []StreamSource
}

type AnimeApi interface {
	GetDetail(SearchResult) Detail
	GetEpisode(SearchResult, uint64) Episode
	Search(string) []SearchResult
}

type AnimixPlayApi struct {
	BaseURL   string
	SearchURL string
}

func NewAnimixPlayApi() AnimixPlayApi {
	return AnimixPlayApi{
		SearchURL: "https://cachecow.eu/api/search",
		BaseURL:   "https://animixplay.to",
	}
}

func (animixApi *AnimixPlayApi) GetEpisode(result SearchResult, number uint64) (Episode, error) {
	res, err := http.Get(result.PageURL)
	if err != nil {
		return Episode{}, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return Episode{}, err

	}

	episodeData := strings.TrimSpace(doc.Find("#epslistplace").Text())
	var data map[string]interface{}
	err = json.Unmarshal([]byte(episodeData), &data)
	if err != nil {
		return Episode{}, err
	}

	urlEp, ok := data[fmt.Sprint(number-1)].(string)
	if !ok {
		return Episode{}, errors.New("url episode was not a string")
	}

	p, err := url.Parse(urlEp)
	if err != nil {
		return Episode{}, err
	}

	streamListHTML, err := http.Get(fmt.Sprintf("https://gogoplay1.com/download?id=%s", p.Query().Get("id")))
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
		streams = append(streams, StreamSource{URL: href, Type: streamType, Origin: "AnimixPlay"})
	})

	return Episode{
		Number:        number,
		StreamSources: streams,
	}, nil
}

func (animixApi *AnimixPlayApi) GetDetail(result SearchResult) (Detail, error) {
	pageHTML, err := http.Get(result.PageURL)
	if err != nil {
		return Detail{}, err
	}

	pageDoc, err := goquery.NewDocumentFromReader(pageHTML.Body)
	if err != nil {
		return Detail{}, err
	}

	animeMetaDataString := pageDoc.Find("script").Last().Text()
	r := regexp.MustCompile("malid = '(.*)'")
	animeID := r.FindStringSubmatch(animeMetaDataString)[1]

	animeMetaDataResp, err := http.Get(fmt.Sprintf("%s/assets/mal/%s.json", animixApi.BaseURL, animeID))
	if err != nil {
		return Detail{}, err
	}

	b, err := io.ReadAll(animeMetaDataResp.Body)
	if err != nil {
		return Detail{}, err
	}

	var animeMetaData map[string]interface{}
	err = json.Unmarshal(b, &animeMetaData)
	if err != nil {
		return Detail{}, err
	}

	return Detail{
		Synopsis: fmt.Sprint(animeMetaData["synopsis"]),
		Type:     fmt.Sprint(animeMetaData["type"]),
		MalURL:   fmt.Sprint(animeMetaData["url"]),
		Episodes: uint64(animeMetaData["episodes"].(float64)),
	}, nil
}

func (animixApi *AnimixPlayApi) Search(name string) ([]SearchResult, error) {
	body := strings.NewReader(fmt.Sprintf("qfast=%s&root=animixplay.to", name))
	resp, err := http.Post(animixApi.SearchURL, "application/x-www-form-urlencoded", body)
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

	htmlString, ok := data["result"].(string)
	if !ok {
		return []SearchResult{}, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf("<html>%s</html>", htmlString)))
	if err != nil {
		return []SearchResult{}, err
	}

	var results []SearchResult
	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		name := s.Find("p.name").Text()
		description := s.Find("p.infotext").Text()
		imgSrc, _ := s.Find("img").Attr("src")
		pageURL, _ := s.Find("a").Attr("href")

		result := SearchResult{
			Title:       name,
			Description: description,
			ImageSrc:    imgSrc,
			PageURL:     animixApi.BaseURL + pageURL,
		}

		results = append(results, result)
	})

	return results, nil
}
