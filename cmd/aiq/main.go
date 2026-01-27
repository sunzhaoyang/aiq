package main

import (
	"flag"
	"os"

	"github.com/aiq/aiq/internal/cli"
)

func main() {
	var sessionFile string
	flag.StringVar(&sessionFile, "s", "", "Path to session file to restore")
	flag.StringVar(&sessionFile, "session", "", "Path to session file to restore")
	flag.Parse()

	if err := cli.Run(sessionFile); err != nil {
		os.Exit(1)
	}
}
