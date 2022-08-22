package cmd

import (
	"okp4/cosmos-faucet/pkg/client"
	"time"

	"github.com/spf13/cobra"
)

// NewSendCommand returns a CLI command to interactively send amount token(s) to given address.
func NewSendCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send <address>",
		Short: "Send tokens to a given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			send := make(chan *client.TriggerTx)
			faucet, err := client.NewFaucet(config, send)
			if err != nil {
				return err
			}

			defer func(faucet *client.Faucet) {
				_ = faucet.Close()
			}(faucet)

			if err := faucet.Send(args[0]); err != nil {
				return err
			}

			send <- client.MakeTriggerTx(client.WithDeadline(time.Now().Add(config.TxTimeout)))

			return err
		},
	}

	return sendCmd
}

func init() {
	rootCmd.AddCommand(NewSendCommand())
}
