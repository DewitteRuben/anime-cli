package main

import (
	"anime-cli/api"
	"anime-cli/cli"
	"anime-cli/video"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	cliArgs := cli.GetCliArgs()

	searchInput, err := cli.PromptSearchAnime()
	if err != nil {
		panic(err)
	}

	api := api.NewAnimixPlayApi()
	results, err := api.Search(searchInput)
	if err != nil {
		panic(err)
	}

	selectedAnime, err := cli.PromptSelectAnime(results)
	if err != nil {
		panic(err)
	}

	animeDetail, err := api.GetDetail(selectedAnime)
	if err != nil {
		panic(err)
	}

	selectedEpisode, err := cli.PromptEpisodeNumber(animeDetail.Episodes)
	if err != nil {
		panic(err)
	}

	ep, err := api.GetEpisode(selectedAnime, selectedEpisode)
	if err != nil {
		panic(err)
	}

	source, err := cli.PromptSelectSource(ep.StreamSources)
	if err != nil {
		panic(err)
	}

	go func() {
		player := video.NewPlayer(cliArgs.Player)
		player.Play(source)
	}()

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		done <- true
	}()

	<-done
}
