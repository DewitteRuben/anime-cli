package cli

import "flag"

type CliArgs struct {
	Player  string
	Verbose bool
}

var cliArgs CliArgs

func InitCliArgs() {
	videoPlayer := flag.String("player", "vlc", "Video player (Supported: VLC, MPV)")
	verbose := flag.Bool("v", false, "Verbose logs")

	flag.Parse()

	cliArgs = CliArgs{
		Player:  *videoPlayer,
		Verbose: *verbose,
	}
}

func GetCliArgs() CliArgs {
	if !flag.Parsed() {
		InitCliArgs()
	}
	return cliArgs
}
