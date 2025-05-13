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

func Execute() error {
	if len(os.Args) > 1 && os.Args[1][0] != '-' {
		if isCommand(os.Args[1]) {
			return rootCmd.Execute()
		}

		codeplugFile = os.Args[1]

		if len(os.Args) > 2 {
			newArgs := make([]string, 0, len(os.Args)-1)
			newArgs = append(newArgs, os.Args[0])
			newArgs = append(newArgs, os.Args[2:]...)
			os.Args = newArgs
		}
	}

	return rootCmd.Execute()
}

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
