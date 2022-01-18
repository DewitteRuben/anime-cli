package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	b64 "encoding/base64"

	"github.com/PuerkitoBio/goquery"
	"github.com/spacemonkeygo/openssl"
)

type GoGoAnimeApi struct {
	BaseURL   string
	SearchURL string
}

type SourceResponse struct {
	Source []struct {
		File  string `json:"file"`
		Label string `json:"label"`
		Type  string `json:"type"`
	} `json:"source"`
	SourceBk []struct {
		File  string `json:"file"`
		Label string `json:"label"`
		Type  string `json:"type"`
	} `json:"source_bk"`
	Track struct {
		Tracks []struct {
			File string `json:"file"`
			Kind string `json:"kind"`
		} `json:"tracks"`
	} `json:"track"`
	Advertising []interface{} `json:"advertising"`
	Linkiframe  string        `json:"linkiframe"`
}

func GetGoGoAnimeSourceData(epId string) (SourceResponse, error) {
	secretKey, err := hex.DecodeString("3235373436353338353932393338333936373634363632383739383333323838")
	if err != nil {
		return SourceResponse{}, err
	}

	iv, err := hex.DecodeString("34323036393133333738303038313335")
	if err != nil {
		return SourceResponse{}, err
	}

	cipher, err := openssl.GetCipherByName("aes-256-cbc")
	if err != nil {
		return SourceResponse{}, err
	}

	ctx, err := openssl.NewEncryptionCipherCtx(cipher, nil, secretKey, iv)
	if err != nil {
		return SourceResponse{}, err
	}

	cipherbytes, err := ctx.EncryptUpdate([]byte(epId))
	if err != nil {
		return SourceResponse{}, err
	}

	finalbytes, err := ctx.EncryptFinal()
	if err != nil {
		return SourceResponse{}, err
	}

	cipherbytes = append(cipherbytes, finalbytes...)

	encryptedId := b64.StdEncoding.EncodeToString(cipherbytes)

	params := url.Values{}
	params.Add("id", encryptedId)
	params.Add("time", `69420691337800813569`)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", os.ExpandEnv("https://gogoplay.io/encrypt-ajax.php"), body)
	if err != nil {
		return SourceResponse{}, err
	}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SourceResponse{}, err
	}

	defer resp.Body.Close()

	respJSON := SourceResponse{}
	json.NewDecoder(resp.Body).Decode(&respJSON)

	return respJSON, nil
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
	r := regexp.MustCompile(`id=([^&]*)`)
	epId := r.FindStringSubmatch(streamListURL)[0][3:]

	sourceData, err := GetGoGoAnimeSourceData(epId)
	if err != nil {
		return Episode{}, err
	}

	var streams []StreamSource
	for _, source := range sourceData.Source {
		streams = append(streams, StreamSource{URL: source.File, Type: source.Label, Origin: gogoApi.Tag()})
	}

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
