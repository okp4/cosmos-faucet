package cmd

import (
	"context"

	"okp4/cosmos-faucet/pkg/client"

	"github.com/spf13/cobra"
)

// NewSendCommand returns a CLI command to interactively send amount token(s) to given address.
func NewSendCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send <address>",
		Short: "Send tokens to a given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			faucet, err := client.NewFaucet(config)
			if err != nil {
				return err
			}

			defer func(faucet *client.Faucet) {
				_ = faucet.Close()
			}(faucet)

			_, err = faucet.SendTxMsg(context.Background(), args[0])

			return err
		},
	}

	return sendCmd
}

func init() {
	rootCmd.AddCommand(NewSendCommand())
}
