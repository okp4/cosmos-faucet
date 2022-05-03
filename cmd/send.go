package cmd

import (
	"okp4/cosmos-faucet/pkg/send"

	"github.com/spf13/cobra"
)

// NewSendCommand returns a CLI command to interactively send amount token(s) to given address.
func NewSendCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send <address>",
		Short: "Send tokens to a given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return send.Send(config, args[0])
		},
	}

	return sendCmd
}

func init() {
	rootCmd.AddCommand(NewSendCommand())
}
