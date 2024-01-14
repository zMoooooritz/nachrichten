
nachrichten
======

[![Latest Release](https://img.shields.io/github/release/zMoooooritz/nachrichten.svg?style=for-the-badge)](https://github.com/zMoooooritz/nachrichten/releases)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/zMoooooritz/nachrichten)
[![Software License](https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)](/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/zMoooooritz/nachrichten/build.yml?branch=master&style=for-the-badge)](https://github.com/zMoooooritz/nachrichten/actions)
[![Go ReportCard](https://goreportcard.com/badge/github.com/zMoooooritz/nachrichten?style=for-the-badge)](https://goreportcard.com/report/zMoooooritz/nachrichten)

Stay informed without leaving your command line

Access up-to-date news in German provided by the [tagesschau](https://www.tagesschau.de/)

<img alt="Welcome to nachrichten" src="https://github.com/zMoooooritz/nachrichten/blob/media/media/demo.gif" width="800" />

## ⇁ Installation 
### Package Manager
```bash
# Arch Linux (btw)
yay -S nachrichten # only AUR
```

### Go
Install directly using `go`
```bash
go install github.com/zMoooooritz/nachrichten@latest
```
or download from [releases](https://github.com/zMoooooritz/nachrichten/releases)

## ⇁ Usage
Run the `nachrichten` command to launch the minimalistic yet informative terminal interface

```bash
Usage of nachrichten:
  -config string
        Path to configuration file
  -debug string
        Path to log file
  -shortnews
    	Only open the current short news
  -version
    	Display version
```

## ⇁ Configuration
The tool does allow for user customization
1. **Theme** - Adapt the used colors in order to change the look and feel of the application
2. **Keybinds** - Customize all keys used within the application
3. **Applications** - Some news related resources can't be shown in a TUI, configure the apps used to open those resources
4. **Settings** - General settings that alter the behavior of the application

An example configuration can be found [here](https://github.com/zMoooooritz/nachrichten/blob/master/configs/config.yaml)

The default keybinds are as follows:

| Key              | Action                 |
| ---------------- | ---------------------- |
| arrows / hjkl    | navigation             |
| g / G            | goto start / end       |
| tab / shift+tab  | change tabs            |
| f                | maximize reader/viewer |
| i                | toggle thumbnail       |
| o                | open article url       |
| v                | open article vod       |
| s                | open current news vod  |
| ?                | toggle help            |
| q / esc / ctrl+c | quit                   |

## ⇁ Built with
- [bubbletea](https://github.com/charmbracelet/bubbletea) and its awesome ecosystem

