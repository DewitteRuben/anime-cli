package main

import (
	"anime-cli/api"
	"anime-cli/cli"
	"anime-cli/storage"
	"anime-cli/video"
	"fmt"
)

func main() {
	cliArgs, err := cli.GetCliArgs()
	if err != nil {
		return
	}

	err = storage.Init()
	if err != nil {
		if cliArgs.Verbose {
			fmt.Println(err)
		}
	}

	animeApi := api.NewApi(cliArgs.AnimeApi)

	for {
		var searchResults []api.SearchResult

		for {
			searchInput, err := cli.PromptSearchAnime()
			if err != nil {
				if cliArgs.Verbose {
					fmt.Println(err)
				}

				return
			}

			searchResults, err = animeApi.Search(searchInput)
			if err != nil {
				if cliArgs.Verbose {
					fmt.Println(err)
				}

				return
			}

			if len(searchResults) > 0 {
				break
			} else {
				fmt.Println("No anime were found for input:", searchInput)
			}
		}

		selectedAnime, err := cli.PromptSelectAnime(searchResults)
		if err != nil {
			if cliArgs.Verbose {
				fmt.Println(err)
			}

			return
		}

		animeDetail, err := animeApi.GetDetail(selectedAnime)
		if err != nil {
			if cliArgs.Verbose {
				fmt.Println(err)
			}

			return
		}

		for {
			selectedEpisode, err := cli.PromptEpisodeNumber(animeDetail.Episodes)
			if err != nil {
				if cliArgs.Verbose {
					fmt.Println(err)
				}

				break
			}

			ep, err := animeApi.GetEpisode(selectedAnime, selectedEpisode)
			if err != nil {
				if cliArgs.Verbose {
					fmt.Println(err)
				}

				break
			}

			source, err := cli.PromptSelectSource(ep.StreamSources)
			if err != nil {
				if cliArgs.Verbose {
					fmt.Println(err)
				}

				break
			}

			prefs := storage.UserPrefs{
				PrefferedApi:      animeApi.Tag(),
				PreferredSource:   source,
				CurrentlyWatching: selectedAnime,
			}

			storage.Persist(prefs)

			player := video.NewPlayer(cliArgs.Player)
			err = player.Play(source)
			if err != nil {
				if cliArgs.Verbose {
					fmt.Println(err)
				}

			}
		}
	}
}
