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
	configuration config.Configuration
}

func NewOpener(configuration config.Configuration) Opener {
	return Opener{
		configuration: configuration,
	}
}

func (o Opener) OpenUrl(t ResourceType, url string) {
	var appConfig config.ApplicationConfig

	switch t {
	case TypeImage:
		appConfig = o.configuration.AppConfig.Image
	case TypeAudio:
		appConfig = o.configuration.AppConfig.Audio
	case TypeVideo:
		appConfig = o.configuration.AppConfig.Video
	case TypeHTML:
		appConfig = o.configuration.AppConfig.HTML
	default:
		defaultOpenUrl(url)
		return
	}

	cConfig := appConfig
	cConfig.Args = append([]string(nil), appConfig.Args...)

	if cConfig.Path == "" || len(cConfig.Args) == 0 {
		defaultOpenUrl(url)
		return
	}

	for i, arg := range cConfig.Args {
		if arg == "$" {
			cConfig.Args[i] = url
		}
	}
	_ = exec.Command(cConfig.Path, cConfig.Args...).Start()
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
