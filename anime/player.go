package anime

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Player interface {
	Play(Stream)
}

type DefaultPlayer struct {
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (dp DefaultPlayer) Play(stream Stream) {
	arguments := []string{stream.URL}
	if stream.Origin == "AnimixPlay" {
		arguments = append(arguments, "--http-referrer='https://gogoplay1.com/'")
	}

	cmd := exec.Command("vlc", arguments...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
}
