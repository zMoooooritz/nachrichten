package util

import (
	"os/exec"
	"runtime"

	"github.com/zMoooooritz/nachrichten/pkg/config"
)

type ResourceType int

const (
	TypeImage ResourceType = iota
	TypeAudio
	TypeVideo
	TypeHTML
)

type Opener struct {
	apps config.Applications
}

func NewOpener(apps config.Applications) Opener {
	return Opener{
		apps: apps,
	}
}

func (o Opener) OpenUrl(t ResourceType, url string) {
	var app config.Application

	switch t {
	case TypeImage:
		app = o.apps.Image
	case TypeAudio:
		app = o.apps.Audio
	case TypeVideo:
		app = o.apps.Video
	case TypeHTML:
		app = o.apps.HTML
	default:
		defaultOpenUrl(url)
		return
	}

	appCopy := app
	appCopy.Args = append([]string(nil), app.Args...)

	if appCopy.Path == "" || len(appCopy.Args) == 0 {
		defaultOpenUrl(url)
		return
	}

	for i, arg := range appCopy.Args {
		if arg == "$" {
			appCopy.Args[i] = url
		}
	}
	_ = exec.Command(appCopy.Path, appCopy.Args...).Start()
}

func defaultOpenUrl(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	_ = exec.Command(cmd, args...).Start()
}
