package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tagesschau"
	"github.com/zMoooooritz/nachrichten/pkg/tui"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

func Run(configFile string, logFile string, shortNews bool) {
	configuration, err := config.Load(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	if logFile != "" {
		err := util.SetLogFile(logFile)
		if err != nil {
			log.Fatalln("Error occoured while setting up the logger: ", err)
		}
		util.Logger.Println("Application started.")
	}

	if shortNews {
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
		// tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
