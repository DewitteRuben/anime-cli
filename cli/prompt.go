package cli

import (
	"anime-cli/api"
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func PromptSelectAnime(searchResults []api.SearchResult) (api.SearchResult, error) {
	if len(searchResults) == 0 {
		return api.SearchResult{}, errors.New("search results are empty")
	}

	selectPrompt := promptui.Select{
		Label: "Select an anime",
		Items: searchResults,
		Templates: &promptui.SelectTemplates{
			Active:   fmt.Sprintf("%s {{ .Title | underline | cyan }}", promptui.IconSelect),
			Inactive: "  {{.Title}}",
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Title | faint }}`, promptui.IconGood),
		},
		Size: 5,
	}

	index, _, err := selectPrompt.Run()
	if err != nil {
		print(err)
		return api.SearchResult{}, err
	}

	return searchResults[index], nil
}

func PromptSelectSource(sources []api.StreamSource) (api.StreamSource, error) {
	if len(sources) == 0 {
		return api.StreamSource{}, errors.New("no streams found")
	}

	selectPrompt := promptui.Select{
		Label: "Select a source",
		Items: sources,
		Templates: &promptui.SelectTemplates{
			Active:   fmt.Sprintf("%s {{ .Type | underline | cyan }}", promptui.IconSelect),
			Inactive: "  {{.Type}}",
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Type | faint }}`, promptui.IconGood),
		},
		Size: 5,
	}

	index, _, err := selectPrompt.Run()
	if err != nil {
		return api.StreamSource{}, err
	}

	return sources[index], nil
}

func PromptEpisodeNumber(episodes uint64) (uint64, error) {
	bold := promptui.Styler(promptui.FGBold)
	searchPrompt := promptui.Prompt{
		Pointer: promptui.PipeCursor,
		Validate: func(s string) error {
			_, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return errors.New("input is not a number")
			}
			return nil
		},
		Templates: &promptui.PromptTemplates{
			Valid: fmt.Sprintf("%s {{ . | bold }}%s ", bold(promptui.IconInitial), bold(":")),
		},
		Label: fmt.Sprintf("Number of episode: [1 - %d]", episodes),
	}

	selectedEpisodeString, err := searchPrompt.Run()
	if err != nil {
		return 0, err
	}

	selectedEpisode, err := strconv.ParseInt(selectedEpisodeString, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint64(selectedEpisode), nil
}

func PromptSearchAnime() (string, error) {
	bold := promptui.Styler(promptui.FGBold)
	searchPrompt := promptui.Prompt{
		Pointer: promptui.PipeCursor,
		Templates: &promptui.PromptTemplates{
			Valid: fmt.Sprintf("%s {{ . | bold }}%s ", bold(promptui.IconInitial), bold(":")),
		},
		Label: "Search for an anime",
	}

	animeSearchInput, err := searchPrompt.Run()
	if err != nil {
		return "", err
	}

	return animeSearchInput, nil
}
