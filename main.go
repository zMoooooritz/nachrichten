package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/tui"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

var (
	// Set via ldflags when building
	Version   = ""
	CommitSHA = ""

	configFile = flag.String("config", "", "Path to configuration file")
	logFile    = flag.String("debug", "", "Path to log file")
	shortNews  = flag.Bool("shortnews", false, "Only open the current short news")
	version    = flag.Bool("version", false, "Display version")
)

func main() {
	flag.Parse()

	if *version {
		if len(CommitSHA) > 7 {
			CommitSHA = CommitSHA[:7]
		}
		if Version == "" {
			Version = "(built from source)"
		}

		fmt.Printf("nachrichten %s", Version)
		if len(CommitSHA) > 0 {
			fmt.Printf(" (%s)", CommitSHA)
		}

		fmt.Println()
		os.Exit(0)
	}

	configuration, err := config.Load(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	if *logFile != "" {
		err := util.SetLogFile(*logFile)
		if err != nil {
			log.Fatalln("Error occoured while setting up the logger: ", err)
		}
		util.Logger.Println("Application started.")
	}

	if *shortNews {
		url, err := tagesschau.GetShortNewsURL()
		if err == nil {
			opener := util.NewOpener(configuration.Applications)
			opener.OpenUrl(util.TypeVideo, url)
		} else {
			log.Fatalln("Error occoured while fetching shornews URL")
		}
		os.Exit(0)
	}

	p := tea.NewProgram(tui.InitialModel(configuration),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
