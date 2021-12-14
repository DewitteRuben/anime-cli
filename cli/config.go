package cli

import "flag"

type CliArgs struct {
	Player   string
	AnimeApi string
	Verbose  bool
}

var cliArgs CliArgs

func InitCliArgs() {
	videoPlayer := flag.String("player", "vlc", "Video player (Supported: VLC, MPV)")
	animeApi := flag.String("api", "animixplay", "Api implementation to fetch videos from (Supported: animixplay)")
	verbose := flag.Bool("v", false, "Verbose logs")

	flag.Parse()

	cliArgs = CliArgs{
		Player:   *videoPlayer,
		AnimeApi: *animeApi,
		Verbose:  *verbose,
	}
}

func GetCliArgs() CliArgs {
	if !flag.Parsed() {
		InitCliArgs()
	}
	return cliArgs
}
