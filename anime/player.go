package anime

import (
	"io"
	"net/http"
	"os"
)

type Player interface {
	Play(Stream)
}

type DefaultPlayer struct {
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (dp DefaultPlayer) Play(stream Stream) {
	go func() {
		httpClient := http.Client{}
		req, err := http.NewRequest("GET", stream.URL, nil)
		if err != nil {
			return
		}

		if stream.Origin == "AnimixPlay" {
			req.Header.Set("referer", "https://gogoplay1.com/")
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return
		}

		defer resp.Body.Close()

		out, err := os.Create("temp")
		if err != nil {
			return
		}
		defer out.Close()

		io.Copy(out, resp.Body)
	}()
}
