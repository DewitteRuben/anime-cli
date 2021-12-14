package api

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
	Tag() AnimeApiTag
	GetDetail(SearchResult) (Detail, error)
	GetEpisode(SearchResult, uint64) (Episode, error)
	Search(string) ([]SearchResult, error)
}

func NewApi(apiTag string) AnimeApi {
	switch apiTag {
	case AnimixPlay.String():
		return NewAnimixPlayApi()
	}
	return nil
}

type AnimeApiTag int64

const (
	AnimixPlay AnimeApiTag = iota
)

func (s AnimeApiTag) String() string {
	switch s {
	case AnimixPlay:
		return "animixplay"
	}
	return "unknown"
}
