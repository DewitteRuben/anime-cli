package storage

import (
	"anime-cli/api"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
)

func Init() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	animeCLIFolder := fmt.Sprintf("%s/%s", homeDir, "/.animecli")
	err = os.MkdirAll(animeCLIFolder, 0777)
	if err != nil {
		return err
	}

	return nil
}

func GetDataFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", homeDir, ".animecli"), nil
}

type UserPrefs struct {
	PrefferedApi      api.AnimeApiTag
	PreferredSource   api.StreamSource
	CurrentlyWatching api.SearchResult
}

func loadPrefs() (UserPrefs, error) {
	dataFolder, err := GetDataFolder()
	if err != nil {
		return UserPrefs{}, err
	}

	filePath := fmt.Sprintf("%s/prefs.json", dataFolder)
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return UserPrefs{}, nil
		}
		return UserPrefs{}, err
	}

	prefs := UserPrefs{}
	err = json.Unmarshal([]byte(file), &prefs)
	if err != nil {
		return UserPrefs{}, err
	}

	return prefs, nil
}

func savePrefs(data UserPrefs) error {
	dataFolder, err := GetDataFolder()
	if err != nil {
		return err
	}

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/prefs.json", dataFolder)
	return ioutil.WriteFile(filePath, file, 0644)
}

func GetPrefs() (UserPrefs, error) {
	return loadPrefs()
}

func Persist(newData UserPrefs) error {
	existingPrefs, err := loadPrefs()
	if err != nil {
		return err
	}

	err = mergo.Merge(&existingPrefs, newData, mergo.WithOverride)
	if err != nil {
		return err
	}

	return savePrefs(existingPrefs)
}
