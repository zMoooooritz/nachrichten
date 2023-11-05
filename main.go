package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	configFile := flag.String("config", "~/.config/nachrichten/config.yaml", "Path to configuration file")

	flag.Parse()

	configuration := loadConfig(*configFile)

	p := tea.NewProgram(InitialModel(configuration),
		tea.WithAltScreen(),
		// tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
