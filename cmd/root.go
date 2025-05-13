package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var codeplugFile string

var rootCmd = &cobra.Command{
	Use:   "anytone-cli",
	Short: "A CLI tool for working with Anytone codeplugs",
	Long: `A command-line interface for working with Anytone codeplugs.
This tool allows you to view and modify parameters in Anytone radio codeplug (.rdt) files
without using the official CPS software.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// If we have arguments and the first one doesn't start with a dash, it's likely our codeplug file
	if len(os.Args) > 1 && os.Args[1][0] != '-' {
		// See if the first arg is a command
		if isCommand(os.Args[1]) {
			// First arg is a command, not a file
			return rootCmd.Execute()
		}

		// First arg is hopefully our codeplug file
		codeplugFile = os.Args[1]

		// Remove the codeplug file from args so cobra doesn't see it
		// This is a bit hacky but works with cobra's arg parsing
		if len(os.Args) > 2 {
			// Create a new slice with everything except the codeplug file
			newArgs := make([]string, 0, len(os.Args)-1)
			newArgs = append(newArgs, os.Args[0])
			newArgs = append(newArgs, os.Args[2:]...)
			os.Args = newArgs
		}
	}

	return rootCmd.Execute()
}

// Check if a string is a known command
func isCommand(cmd string) bool {
	commands := []string{"help", "completion", "info", "set"}
	for _, c := range commands {
		if c == cmd {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(setRadioCmd)
}
