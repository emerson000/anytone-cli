package cmd

import (
	"fmt"
	"strconv"

	"github.com/emerson000/anytone-cli/pkg/codeplug"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get codeplug parameters",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if codeplugFile == "" {
			return fmt.Errorf("codeplug file path is required")
		}
		return nil
	},
}

var getRadioIDCmd = &cobra.Command{
	Use:   "radio_id [index]",
	Short: "Get radio ID(s). If no index is provided, returns all radio IDs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cp, err := codeplug.Open(codeplugFile)
		if err != nil {
			return fmt.Errorf("failed to open codeplug: %w", err)
		}
		defer cp.Close()

		// If no index provided, get all radio IDs
		if len(args) == 0 {
			radioIDs, err := cp.GetRadioIDs()
			if err != nil {
				return fmt.Errorf("failed to get radio IDs: %w", err)
			}
			for _, entry := range radioIDs {
				fmt.Printf("%d: %d (%s)\n", entry.Index, entry.ID, entry.Name)
			}
			return nil
		}

		// Get specific radio ID by index
		index, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}

		radioID, err := cp.GetRadioIDByIndex(index)
		if err != nil {
			return fmt.Errorf("failed to get radio ID: %w", err)
		}
		fmt.Printf("%d: %d (%s)\n", index, radioID.ID, radioID.Name)

		return nil
	},
}

func init() {
	getCmd.AddCommand(getRadioIDCmd)
}
