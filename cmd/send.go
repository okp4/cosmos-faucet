package cmd

import (
	"github.com/spf13/cobra"
)

// SendCommand returns a CLI command to interactively send amount to given address.
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send token to a given address",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
