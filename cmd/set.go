package cmd

import (
	"fmt"
	"strconv"

	"github.com/emerson000/anytone-cli/pkg/codeplug"
	"github.com/spf13/cobra"
)

var setRadioCmd = &cobra.Command{
	Use:   "set",
	Short: "Set codeplug parameters",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if codeplugFile == "" {
			return fmt.Errorf("codeplug file path is required")
		}
		return nil
	},
}

var setRadioIDCmd = &cobra.Command{
	Use:   "radio_id <index> <new_id>",
	Short: "Update a radio ID",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		index, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}

		newID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid radio ID: %w", err)
		}

		cp, err := codeplug.Open(codeplugFile)
		if err != nil {
			return fmt.Errorf("failed to open codeplug: %w", err)
		}
		defer cp.Close()

		if err := cp.UpdateRadioID(index, newID); err != nil {
			return fmt.Errorf("failed to update radio ID: %w", err)
		}

		fmt.Printf("Successfully updated radio ID at index %d to %d\n", index, newID)
		return nil
	},
}

func init() {
	setRadioCmd.AddCommand(setRadioIDCmd)
}
