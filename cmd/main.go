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

func Run() {

	cfg := config.Init()
	configuration := cfg.Load()
	if cfg.LogFile != "" {
		err := util.SetLogFile(cfg.LogFile)
		log.Fatalln("Error occoured while setting up logger: ", err)
	}
	util.Logger.Println("Application started.")

	p := tea.NewProgram(tui.InitialModel(configuration),
		tea.WithAltScreen(),
		// tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
