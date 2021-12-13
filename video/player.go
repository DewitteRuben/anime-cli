package video

import (
	"anime-cli/api"
	"os/exec"
)

type Player interface {
	Play(api.StreamSource)
}

func NewVLCPlayer() VLC {
	return VLC{}
}

func NewMPVPlayer() MPV {
	return MPV{}
}

func NewPlayer(playerTag string) Player {
	switch playerTag {
	case "mpv":
		return NewMPVPlayer()
	case "vlc":
		return NewVLCPlayer()
	}
	return nil
}

type MPV struct{}
type VLC struct{}

func (mpv MPV) Play(stream api.StreamSource) {
	arguments := []string{stream.URL}
	if stream.Origin == "AnimixPlay" {
		arguments = append(arguments, "--http-header-fields='referrer: https://gogoplay1.com/'")
	}

	_ = runCommand("mpv", arguments)
}

func runCommand(command string, arguments []string) error {
	return exec.Command(command, arguments...).Run()
}

func (dp VLC) Play(stream api.StreamSource) {
	arguments := []string{stream.URL}
	if stream.Origin == "AnimixPlay" {
		arguments = append(arguments, "--http-referrer='https://gogoplay1.com/'")
	}

	_ = runCommand("vlc", arguments)
}
