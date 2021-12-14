package main

import (
	"anime-cli/api"
	"anime-cli/cli"
	"anime-cli/video"
	"log"
)

func main() {
	cliArgs, err := cli.GetCliArgs()
	if err != nil {
		return
	}

	for {
		searchInput, err := cli.PromptSearchAnime()
		if err != nil {
			if cliArgs.Verbose {
				log.Println(err)
			}

			return
		}

		api := api.NewApi(cliArgs.AnimeApi)
		results, err := api.Search(searchInput)
		if err != nil {
			if cliArgs.Verbose {
				log.Println(err)
			}

			return
		}

		selectedAnime, err := cli.PromptSelectAnime(results)
		if err != nil {
			if cliArgs.Verbose {
				log.Println(err)
			}

			return
		}

		animeDetail, err := api.GetDetail(selectedAnime)
		if err != nil {
			if cliArgs.Verbose {
				log.Println(err)
			}

			return
		}

		for {
			selectedEpisode, err := cli.PromptEpisodeNumber(animeDetail.Episodes)
			if err != nil {
				if cliArgs.Verbose {
					log.Println(err)
				}

				break
			}

			ep, err := api.GetEpisode(selectedAnime, selectedEpisode)
			if err != nil {
				if cliArgs.Verbose {
					log.Println(err)
				}

				break
			}

			source, err := cli.PromptSelectSource(ep.StreamSources)
			if err != nil {
				if cliArgs.Verbose {
					log.Println(err)
				}

				break
			}

			player := video.NewPlayer(cliArgs.Player)
			player.Play(source)
		}
	}
}
