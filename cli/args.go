package cli

import (
	"flag"
	"os"

	"github.com/jessevdk/go-flags"
)

type CliArgs struct {
	Player   string `short:"p" long:"player" description:"Video player to play videos with" choice:"vlc" choice:"mpv"`
	AnimeApi string `short:"a" long:"api" description:"Site to fetch data and stream URLs from" choice:"gogoanime"`
	Verbose  bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
}

var cliArgs = CliArgs{
	Player:   "vlc",
	AnimeApi: "gogoanime",
	Verbose:  false,
}

func InitCliArgs() ([]string, error) {
	return flags.ParseArgs(&cliArgs, os.Args)
}

func GetCliArgs() (CliArgs, error) {
	if !flag.Parsed() {
		_, err := InitCliArgs()
		if err != nil {
			return CliArgs{}, err
		}
	}
	return cliArgs, nil
}
