package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zMoooooritz/nachrichten/cmd"
)

var (
	// Set via ldflags when building
	Version   = ""
	CommitSHA = ""

	configFile = flag.String("config", "", "Path to configuration file")
	logFile    = flag.String("debug", "", "Path to log file")
	shortnews  = flag.Bool("shortnews", false, "Only open the current short news")
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

	cmd.Run(*configFile, *logFile, *shortnews)
}
