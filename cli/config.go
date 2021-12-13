package cli

import "flag"

type CliArgs struct {
	Player string
}

func GetCliArgs() CliArgs {
	videoPlayer := flag.String("player", "vlc", "Video player (Supported: VLC, MPV)")
	flag.Parse()

	return CliArgs{
		Player: *videoPlayer,
	}
}
