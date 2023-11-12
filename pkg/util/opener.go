package util

import (
	"os/exec"
	"runtime"

	"github.com/zMoooooritz/nachrichten/pkg/config"
)

func OpenUrl(t config.ResourceType, c config.Configuration, url string) error {
	var appConfig config.ApplicationConfig

	switch t {
	case config.TypeImage:
		appConfig = c.AppConfig.Image
	case config.TypeAudio:
		appConfig = c.AppConfig.Audio
	case config.TypeVideo:
		appConfig = c.AppConfig.Video
	case config.TypeHTML:
		appConfig = c.AppConfig.HTML
	default:
		return defaultOpenUrl(url)
	}

	cConfig := appConfig
	cConfig.Args = append([]string(nil), appConfig.Args...)

	if cConfig.Path == "" || len(cConfig.Args) == 0 {
		return defaultOpenUrl(url)
	}

	for i, arg := range cConfig.Args {
		if arg == "$" {
			cConfig.Args[i] = url
		}
	}
	return exec.Command(cConfig.Path, cConfig.Args...).Start()
}

func defaultOpenUrl(url string) error {
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
	return exec.Command(cmd, args...).Start()
}
