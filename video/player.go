package video

import (
	"anime-cli/api"
	"anime-cli/cli"
	"os"
	"os/exec"
)

type Player interface {
	Play(api.StreamSource) error
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

func (mpv MPV) Play(stream api.StreamSource) error {
	arguments := []string{}
	if stream.Origin == "AnimixPlay" {
		arguments = append(arguments, "--http-header-fields=Referer: https://gogoplay1.com/")
	}
	arguments = append(arguments, stream.URL)

	return runCommand("mpv", arguments)
}

func runCommand(command string, arguments []string) error {
	cmd := exec.Command(command, arguments...)
	config := cli.GetCliArgs()
	if config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

func (dp VLC) Play(stream api.StreamSource) error {
	arguments := []string{stream.URL}
	if stream.Origin == "AnimixPlay" {
		arguments = append(arguments, "--http-referrer='https://gogoplay1.com/'")
	}

	return runCommand("vlc", arguments)
}
