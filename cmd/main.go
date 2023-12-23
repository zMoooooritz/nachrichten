package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zMoooooritz/nachrichten/pkg/config"
	"github.com/zMoooooritz/nachrichten/pkg/tui"
	"github.com/zMoooooritz/nachrichten/pkg/util"
)

func Run(configFile string, logFile string) {
	configuration := config.Load(configFile)
	err := util.SetLogFile(logFile)
	if err != nil {
		log.Fatalln("Error occoured while setting up logger: ", err)
	}
	util.Logger.Println("Application started.")

	p := tea.NewProgram(tui.InitialModel(configuration),
		tea.WithAltScreen(),
		// tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
