package main

import (
	"fmt"
	"os"

	"github.com/aiq/aiq/internal/cli"
	"github.com/spf13/cobra"
)

var sessionFile string

var rootCmd = &cobra.Command{
	Use:   "aiq",
	Short: "AIQ - An intelligent SQL client",
	Long: `AIQ (AI Query) is an intelligent SQL client that translates your 
natural language questions into precise SQL queries for MySQL and other databases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.Run(sessionFile)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&sessionFile, "session", "s", "", "Restore conversation from session file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
