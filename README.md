
# anime-cli

Watch your favourite anime using the video player of your choice directly from the command line.


## Build

Clone this repository and run:

```
go mod download
go build .
```

## Install

Download the latest pre-compiled binary that fits your OS and architecture [here](https://github.com/DewitteRuben/anime-cli/releases/latest)

If your architecture and/or your OS is not in the list, please consider compiling your own using Golang.

## Usage

```
anime-cli [OPTIONS]

Application Options:
  -p, --player=[vlc|mpv]           Video player to play videos with (default: vlc)
  -a, --api=[animixplay|gogoanime] Site to fetch data and stream URLs from (default: gogoanime)
  -v, --verbose                    Show verbose debug information

Help Options:
  -h, --help                       Show this help message
```

## Dependencies

* vlc or mpv (make sure any or both of them are exposed in $PATH)
