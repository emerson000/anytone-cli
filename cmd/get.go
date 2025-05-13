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

var getChannelCmd = &cobra.Command{
	Use:   "channel [index]",
	Short: "Get channel(s). If no index is provided, returns all channels.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cp, err := codeplug.Open(codeplugFile)
		if err != nil {
			return fmt.Errorf("failed to open codeplug: %w", err)
		}
		defer cp.Close()

		if len(args) == 0 {
			channels, err := cp.GetChannels()
			if err != nil {
				return fmt.Errorf("failed to get channels: %w", err)
			}
			for i, channel := range channels {
				fmt.Printf("%d: %s (Rx: %.4f MHz, Tx: %.4f MHz)\n", i, channel.Name, float64(channel.RxFreq)/100000, float64(channel.TxFreq)/100000)
			}
			return nil
		}

		index, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}

		channel, err := cp.GetChannelByIndex(index)
		if err != nil {
			return fmt.Errorf("failed to get channel: %w", err)
		}

		fmt.Printf("Channel %d:\n", index)
		fmt.Printf("  Name: %s\n", channel.Name)
		fmt.Printf("  Rx Frequency: %.4f MHz\n", float64(channel.RxFreq)/100000)
		fmt.Printf("  Tx Frequency: %.4f MHz\n", float64(channel.TxFreq)/100000)
		fmt.Printf("  Channel Type: %d\n", channel.ChannelType)
		fmt.Printf("  Tx Power: %d\n", channel.TxPower)
		fmt.Printf("  Bandwidth: %d\n", channel.Bandwidth)
		fmt.Printf("  CTCSS/DCS Decode: %d\n", channel.CtcssDcsDecode)
		fmt.Printf("  CTCSS/DCS Encode: %d\n", channel.CtcssDcsEncode)
		fmt.Printf("  Radio ID: %d\n", channel.RadioId)
		fmt.Printf("  Scan List: %d\n", channel.ScanList)
		fmt.Printf("  Color Code: %d\n", channel.RxColorCode)
		fmt.Printf("  Slot: %d\n", channel.Slot)

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
	getCmd.AddCommand(getChannelCmd)
}
