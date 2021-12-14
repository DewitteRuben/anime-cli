package main

import (
	"anime-cli/api"
	"anime-cli/cli"
	"anime-cli/video"
	"log"
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
		if cliArgs.Verbose {
			log.Fatalln(err)
		}
	}

	api := api.NewAnimixPlayApi()
	results, err := api.Search(searchInput)
	if err != nil {
		if cliArgs.Verbose {
			log.Fatalln(err)
		}
	}

	selectedAnime, err := cli.PromptSelectAnime(results)
	if err != nil {
		if cliArgs.Verbose {
			log.Fatalln(err)
		}
	}

	animeDetail, err := api.GetDetail(selectedAnime)
	if err != nil {
		if cliArgs.Verbose {
			log.Fatalln(err)
		}
	}

	go func() {
		for {
			selectedEpisode, err := cli.PromptEpisodeNumber(animeDetail.Episodes)
			if err != nil {
				if cliArgs.Verbose {
					log.Fatalln(err)
				}
			}

			ep, err := api.GetEpisode(selectedAnime, selectedEpisode)
			if err != nil {
				if cliArgs.Verbose {
					log.Fatalln(err)
				}
			}

			source, err := cli.PromptSelectSource(ep.StreamSources)
			if err != nil {
				if cliArgs.Verbose {
					log.Fatalln(err)
				}
			}

			player := video.NewPlayer(cliArgs.Player)
			player.Play(source)
		}
	}()

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		done <- true
	}()

	<-done
}
