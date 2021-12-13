package main

import (
	"anime-cli/anime"
)

func main() {
	api := anime.NewAnimixPlayApi()
	results, _ := api.Search("death note")

	ep, _ := api.GetEpisode(results[0], 1)

	player := anime.DefaultPlayer{}
	player.Play(ep.Streams[3])

}
