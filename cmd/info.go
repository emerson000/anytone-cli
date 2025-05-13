package cmd

import (
	"fmt"

	"github.com/emerson000/anytone-cli/pkg/codeplug"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display information about the codeplug",
	RunE: func(cmd *cobra.Command, args []string) error {
		if codeplugFile == "" {
			return fmt.Errorf("codeplug file path is required")
		}

		cp, err := codeplug.Open(codeplugFile)
		if err != nil {
			return fmt.Errorf("failed to open codeplug: %w", err)
		}
		defer cp.Close()

		info, err := cp.GetInfo()
		if err != nil {
			return fmt.Errorf("failed to get codeplug info: %w", err)
		}

		fmt.Printf("Model: %s\n", info.Model)
		fmt.Printf("Radio IDs:\n")
		for i, id := range info.RadioIDs {
			fmt.Printf("  %d: %d\n", info.RadioIDIndices[i], id)
		}

		return nil
	},
}
